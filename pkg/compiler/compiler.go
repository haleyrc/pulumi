// Copyright 2016 Marapongo, Inc. All rights reserved.

package compiler

import (
	"os"

	"github.com/golang/glog"
	"github.com/satori/go.uuid"

	"github.com/marapongo/mu/pkg/ast"
	"github.com/marapongo/mu/pkg/compiler/backends"
	"github.com/marapongo/mu/pkg/compiler/backends/clouds"
	"github.com/marapongo/mu/pkg/compiler/backends/schedulers"
	"github.com/marapongo/mu/pkg/compiler/core"
	"github.com/marapongo/mu/pkg/diag"
	"github.com/marapongo/mu/pkg/errors"
	"github.com/marapongo/mu/pkg/workspace"
)

// Compiler provides an interface into the many phases of the Mu compilation process.
type Compiler interface {
	core.Phase

	// Context returns the current compiler context.
	Context() *core.Context

	// Build detects and compiles inputs from the given location, storing build artifacts in the given destination.
	Build(inp string, outp string)
	// BuildFile uses the given Mufile directly, and stores build artifacts in the given destination.
	BuildFile(mufile []byte, ext string, outp string)
}

// compiler is the canonical implementation of the Mu compiler.
type compiler struct {
	ctx  *core.Context
	opts Options
	deps ast.BoundDependencies // a cache of loaded dependencies.
}

// NewCompiler creates a new instance of the Mu compiler, with the given initialization settings.
func NewCompiler(opts Options) Compiler {
	return &compiler{
		ctx:  &core.Context{},
		opts: opts,
		deps: make(ast.BoundDependencies),
	}
}

func (c *compiler) Context() *core.Context {
	return c.ctx
}

func (c *compiler) Diag() diag.Sink {
	return c.opts.Diag
}

func (c *compiler) Build(inp string, outp string) {
	glog.Infof("Building target '%v' (out='%v')", inp, outp)

	// First find the root of the current package based on the location of its Mufile.
	w, err := workspace.New(inp, c.Diag())
	if err != nil {
		c.Diag().Errorf(errors.IOError.WithFile(inp), err)
		return
	}
	mufile, err := w.DetectMufile()
	if err != nil {
		c.Diag().Errorf(errors.IOError.WithFile(inp), err)
		return
	}
	if mufile == "" {
		c.Diag().Errorf(errors.MissingMufile, inp)
		return
	}

	// Read in the contents of the document and make it available to subsequent stages.
	doc, err := diag.ReadDocument(mufile)
	if err != nil {
		c.Diag().Errorf(errors.CouldNotReadMufile.WithFile(mufile), err)
		return
	}

	c.buildDocument(w, doc, outp)
}

func (c *compiler) BuildFile(mufile []byte, ext string, outp string) {
	glog.Infof("Building in-memory %v file (bytes=%v out='%v')", ext, len(mufile), outp)

	// Default to the current working directory for the workspace.
	dir, err := os.Getwd()
	if err != nil {
		c.Diag().Errorf(errors.IOError, err)
		return
	}
	w, err := workspace.New(dir, c.Diag())
	if err != nil {
		c.Diag().Errorf(errors.IOError, err)
		return
	}

	doc := &diag.Document{File: workspace.Mufile + ext, Body: mufile}
	c.buildDocument(w, doc, outp)
}

func (c *compiler) buildDocument(w workspace.W, doc *diag.Document, outp string) {
	glog.Infof("Building doc '%v' (bytes=%v out='%v')", doc.File, len(doc.Body), outp)
	if glog.V(2) {
		defer func() {
			glog.V(2).Infof("Building doc '%v' completed w/ %v warnings and %v errors",
				doc.File, c.Diag().Warnings(), c.Diag().Errors())
		}()
	}

	// Perform the front-end passes to generate a stack AST.
	stack, ok := c.parseStack(doc)
	if !ok {
		return
	}

	// Now expand all dependencies so they are available to semantic analysis.
	deps, ok := c.loadDependencies(w, doc, stack)
	if !ok {
		return
	}

	// Perform semantic analysis on all stacks passes to validate, transform, and/or update the AST.
	stack, ok = c.analyzeStack(doc, stack)
	if !ok {
		return
	}
	for _, dep := range deps {
		if _, ok := c.analyzeStack(doc, dep.Stack); !ok {
			return
		}
	}

	if c.opts.SkipCodegen {
		glog.V(2).Infof("Skipping code-generation (opts.SkipCodegen=true)")
	} else {
		// Figure out which cloud architecture we will be targeting during code-gen.
		target, arch, ok := c.discoverClusterArch(w, doc, stack)
		if !ok {
			return
		}
		if glog.V(2) {
			tname := "n/a"
			if target != nil {
				tname = target.Name
			}
			glog.V(2).Infof("Stack %v targets target=%v cloud=%v", stack.Name, tname, arch)
		}

		// Now get the backend cloud provider to process the stack from here on out.
		be := backends.New(arch, c.Diag())
		be.CodeGen(core.Compiland{target, doc, stack, deps})
	}
}

// parseStack takes a Mufile document, parses and validates it, and returns a stack AST.  If anything goes wrong
// during this process, the number of errors will be non-zero, and the bool will be false.
func (c *compiler) parseStack(doc *diag.Document) (*ast.Stack, bool) {
	// To build the Mu package, first parse the input file.
	p := NewParser(c)
	stack := p.Parse(doc)
	if !p.Diag().Success() {
		// If any errors happened during parsing, exit.
		return stack, false
	}

	// Do a pass over the parse tree to ensure that all is well.
	ptAnalyzer := NewPTAnalyzer(c)
	ptAnalyzer.Analyze(doc, stack)
	if !p.Diag().Success() {
		// If any errors happened during parse tree analysis, exit.
		return stack, false
	}

	return stack, true
}

// loadDependencies enumerates all of the target stack's dependencies, and parses them into AST form.
func (c *compiler) loadDependencies(w workspace.W, doc *diag.Document, stack *ast.Stack) (ast.BoundDependencies, bool) {
	stack.BoundDependencies = make(ast.BoundDependencies)

	ok := true
	for _, ref := range ast.StableDependencies(stack.Dependencies) {
		dep := stack.Dependencies[ref]
		glog.V(3).Infof("Loading Stack %v dependency %v:%v", stack.Name, ref, dep)

		// First see if we've already loaded this dependency.  In that case, we can reuse it.
		var bound *ast.BoundDependency
		if bd, exists := c.deps[ref]; exists {
			bound = &bd
		} else {
			bound = c.loadDependency(w, doc, ref, dep)
			if bound != nil {
				c.deps[ref] = *bound
			}
		}

		if bound == nil {
			// Missing dependency; return false to the caller so we can stop before things get worse.
			ok = false
		} else {
			// TODO: check for version mismatches.
			stack.BoundDependencies[ref] = *bound

			// Now recursively load this stack's dependenciess too.  We won't return them, however, they need to exist
			// on the ASTs so that we can use dependency information during code-generation, for example.
			c.loadDependencies(w, doc, bound.Stack)
		}
	}

	return stack.BoundDependencies, ok
}

// loadDependency loads up the target dependency from the current workspace using the stack resolution rules.
func (c *compiler) loadDependency(w workspace.W, doc *diag.Document, ref ast.Ref,
	dep ast.Dependency) *ast.BoundDependency {
	// There are many places a dependency could come from.  Consult the workspace for a list of those paths.  It will
	// return a number of them, in preferred order, and we simply probe each one until we find something.
	for _, loc := range w.DepCandidates(ref) {
		// Try to read this location as a document.
		isMufile := workspace.IsMufile(loc, c.Diag())
		glog.V(5).Infof("Probing for dependency %v at %v: %v", ref, loc, isMufile)

		if isMufile {
			doc, err := diag.ReadDocument(loc)
			if err != nil {
				c.Diag().Errorf(errors.CouldNotReadMufile.WithFile(loc), err)
				return nil
			}

			// If we got this far, we've loaded up the dependency's Mufile; parse it and return the result.
			// TODO: it's not clear how much "validation" to perform here.  If the target won't compile, for example,
			//     we are going to get weird errors and failure modes.
			stack, ok := c.parseStack(doc)
			if !ok {
				return nil
			}

			return &ast.BoundDependency{
				Ref:     ref,
				Version: ast.SemVer(dep),
				Stack:   stack,
			}
		}
	}

	// If we got to this spot, we could not find the dependency.  Issue an error and bail out.
	c.Diag().Errorf(errors.MissingDependency.WithDocument(doc), ref)
	return nil
}

// analyzeStack performs semantic analysis on a stack -- validating, transforming, and/or updating it -- and then
// returns the result.  If a problem occurs, errors will have been emitted, and the bool return will be false.
func (c *compiler) analyzeStack(doc *diag.Document, stack *ast.Stack) (*ast.Stack, bool) {
	binder := NewBinder(c)
	binder.Bind(doc, stack)
	if !c.Diag().Success() {
		// If any errors happened during binding, exit.
		return stack, false
	}

	// TODO: perform semantic analysis on the bound tree.

	return stack, true
}

// discoverClusterArch uses a variety of mechanisms to discover the target architecture, returning it.  If no
// architecture was discovered, an error is issued, and the bool return will be false.
func (c *compiler) discoverClusterArch(w workspace.W, doc *diag.Document,
	stack *ast.Stack) (*ast.Cluster, backends.Arch, bool) {
	// Cluster and architectures settings may come from one of three places, in order of search preference:
	//		1) command line arguments.
	//		2) settings specific to this stack.
	//		3) cluster-wide settings in a workspace.
	// In other words, 1 overrides 2 which overrides 3.
	arch := c.opts.Arch

	// If a cluster was specified, look it up and load up its options.
	var cluster *ast.Cluster
	if c.opts.Cluster != "" {
		// First, check the stack to see if it has a targets section.
		if cl, exists := stack.Clusters[c.opts.Cluster]; exists {
			cluster = &cl
		} else {
			// If that didn't work, see if the workspace has an opinion.
			if cl, exists := w.Settings().Clusters[c.opts.Cluster]; exists {
				cluster = &cl
			} else {
				c.Diag().Errorf(errors.ClusterNotFound.WithDocument(doc), c.opts.Cluster)
				return cluster, arch, false
			}
		}
	}

	// If no cluster was specified or discovered yet, see if there is a default one to use.
	if cluster == nil {
		for _, cl := range stack.Clusters {
			if cl.Default {
				cluster = &cl
				break
			}
		}
		if cluster == nil {
			for _, cl := range w.Settings().Clusters {
				if cl.Default {
					cluster = &cl
					break
				}
			}
		}
	}

	if cluster == nil {
		// If no target was found, and we don't have an architecture, error out.
		if arch.Cloud == clouds.NoArch {
			c.Diag().Errorf(errors.MissingTarget.WithDocument(doc))
			return cluster, arch, false
		}

		// If we got here, generate an "anonymous" cluster, so that we at least have a name.
		cluster = c.newAnonCluster(arch)
	} else {
		// If a target was found, go ahead and extract and validate the target architecture.
		a, ok := c.getClusterArch(doc, cluster, arch)
		if !ok {
			return cluster, arch, false
		}
		arch = a
	}

	return cluster, arch, true
}

// newAnonCluster creates an anonymous cluster for stacks that didn't declare one.
func (c *compiler) newAnonCluster(arch backends.Arch) *ast.Cluster {
	// TODO: ensure this is unique.
	// TODO: we want to cache names somewhere (~/.mu/?) so that we can reuse temporary local stacks, etc.
	return &ast.Cluster{
		Name:      uuid.NewV4().String(),
		Cloud:     clouds.Names[arch.Cloud],
		Scheduler: schedulers.Names[arch.Scheduler],
	}
}

// getClusterArch gets and validates the architecture from an existing target.
func (c *compiler) getClusterArch(doc *diag.Document, cluster *ast.Cluster,
	existing backends.Arch) (backends.Arch, bool) {
	targetCloud := existing.Cloud
	targetScheduler := existing.Scheduler

	// If specified, look up the cluster's architecture settings.
	if cluster.Cloud != "" {
		tc, ok := clouds.Values[cluster.Cloud]
		if !ok {
			c.Diag().Errorf(errors.UnrecognizedCloudArch.WithDocument(doc), cluster.Cloud)
			return existing, false
		}
		targetCloud = tc
	}
	if cluster.Scheduler != "" {
		ts, ok := schedulers.Values[cluster.Scheduler]
		if !ok {
			c.Diag().Errorf(errors.UnrecognizedSchedulerArch.WithDocument(doc), cluster.Scheduler)
			return existing, false
		}
		targetScheduler = ts
	}

	// Ensure there aren't any conflicts, comparing compiler options to cluster settings.
	tarch := backends.Arch{targetCloud, targetScheduler}
	if targetCloud != existing.Cloud && existing.Cloud != clouds.NoArch {
		c.Diag().Errorf(errors.ConflictingClusterArchSelection.WithDocument(doc), existing, cluster.Name, tarch)
		return tarch, false
	}
	if targetScheduler != existing.Scheduler && existing.Scheduler != schedulers.NoArch {
		c.Diag().Errorf(errors.ConflictingClusterArchSelection.WithDocument(doc), existing, cluster.Name, tarch)
		return tarch, false
	}

	return tarch, true
}
