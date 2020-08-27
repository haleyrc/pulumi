// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
// +build python all

package ints

var dirs = []string{
	"simple",
}

func TestPythonTransformations(t *testing.T) {
	for _, dir := range dirs {
		d := filepath.Join("python", dir)
		t.Run(d, func(t *testing.T) {
			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir: d,
				Dependencies: []string{
					filepath.Join("..", "..", "..", "sdk", "python", "env", "src"),
				},
				Quick:                  true,
				ExtraRuntimeValidation: validator("python"),
			})
		})
	}
}
