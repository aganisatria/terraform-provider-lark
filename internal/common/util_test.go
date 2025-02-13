// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"strings"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		substrs []string
		want    bool
		mock    bool
	}{
		{
			name:    "string empty with empty substring",
			s:       "",
			substrs: []string{},
			want:    false,
			mock:    false,
		},
		{
			name:    "string empty with substring",
			s:       "",
			substrs: []string{"test"},
			want:    false,
			mock:    false,
		},
		{
			name:    "string with one substring match",
			s:       "hello world",
			substrs: []string{"world"},
			want:    true,
			mock:    false,
		},
		{
			name:    "string with multiple substrings, one match",
			s:       "hello world",
			substrs: []string{"foo", "world", "bar"},
			want:    true,
			mock:    false,
		},
		{
			name:    "string with multiple substrings, no match",
			s:       "hello world",
			substrs: []string{"foo", "bar", "baz"},
			want:    false,
			mock:    false,
		},
		{
			name:    "test with mock strings.Contains always false",
			s:       "hello world",
			substrs: []string{"world"},
			want:    false,
			mock:    true,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			if tt.mock {
				Mock(strings.Contains).To(func(s, substr string) bool {
					return false
				}).Build()
			}

			got := contains(tt.s, tt.substrs...)
			So(got, ShouldEqual, tt.want)
		})
	}
}
