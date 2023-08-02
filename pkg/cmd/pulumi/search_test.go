// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/backend"
	"github.com/pulumi/pulumi/pkg/v3/backend/display"
	"github.com/pulumi/pulumi/pkg/v3/backend/httpstate"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_cmd(t *testing.T) {
	t.Parallel()
	var buff bytes.Buffer
	name := "foo"
	typ := "bar"
	program := "program1"
	stack := "stack1"
	pack := "pack1"
	mod := "mod1"
	modified := "2023-01-01T00:00:00.000Z"
	cmd := searchCmd{
		Stdout: &buff,
		currentBackend: func(context.Context, *workspace.Project, display.Options) (backend.Backend, error) {
			return &stubHTTPBackend{
				SearchF: func(context.Context, string, *apitype.PulumiQueryRequest) (*apitype.ResourceSearchResponse, error) {
					return &apitype.ResourceSearchResponse{
						Resources: []apitype.ResourceResult{
							{
								Name:     &name,
								Type:     &typ,
								Program:  &program,
								Stack:    &stack,
								Package:  &pack,
								Module:   &mod,
								Modified: &modified,
							},
						},
					}, nil
				},
				CurrentUserF: func() (string, []string, error) {
					return "user", []string{"org1", "org2"}, nil
				},
			}, nil
		},
	}

	err := cmd.Run(context.Background(), []string{})
	require.NoError(t, err)

	assert.Contains(t, buff.String(), name)
	assert.Contains(t, buff.String(), typ)
	assert.Contains(t, buff.String(), program)
}

type stubHTTPBackend struct {
	httpstate.Backend

	SearchF func(
		context.Context, string, *apitype.PulumiQueryRequest,
	) (*apitype.ResourceSearchResponse, error)
	NaturalLanguageSearchF func(context.Context, string, string) (*apitype.ResourceSearchResponse, error)
	CurrentUserF           func() (string, []string, error)
}

var _ httpstate.Backend = (*stubHTTPBackend)(nil)

func (f *stubHTTPBackend) Search(
	ctx context.Context, orgName string, queryParams *apitype.PulumiQueryRequest,
) (*apitype.ResourceSearchResponse, error) {
	return f.SearchF(ctx, orgName, queryParams)
}

func (f *stubHTTPBackend) CurrentUser() (string, []string, error) {
	return f.CurrentUserF()
}
