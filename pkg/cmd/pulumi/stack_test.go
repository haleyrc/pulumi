// Copyright 2016-2013, Pulumi Corporation.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringifyOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give any
		want string
	}{
		{"int", 42, "42"},
		{"string", "ABC", "ABC"},
		{
			desc: "array",
			give: []string{"hello", "goodbye"},
			want: `["hello","goodbye"]`,
		},
		{
			desc: "object",
			give: map[string]any{
				"foo": 42,
				"bar": map[string]any{
					"baz": true,
				},
			},
			want: `{"bar":{"baz":true},"foo":42}`,
		},
		{
			desc: "special characters",
			give: "pass&word",
			want: "pass&word",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			got := stringifyOutput(tt.give)
			assert.Equal(t, tt.want, got)
		})
	}
}
