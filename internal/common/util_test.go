package common

import (
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
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
			name:    "string kosong tanpa substring",
			s:       "",
			substrs: []string{},
			want:    false,
			mock:    false,
		},
		{
			name:    "string kosong dengan substring",
			s:       "",
			substrs: []string{"test"},
			want:    false,
			mock:    false,
		},
		{
			name:    "string dengan satu substring yang cocok",
			s:       "hello world",
			substrs: []string{"world"},
			want:    true,
			mock:    false,
		},
		{
			name:    "string dengan beberapa substring, satu cocok",
			s:       "hello world",
			substrs: []string{"foo", "world", "bar"},
			want:    true,
			mock:    false,
		},
		{
			name:    "string dengan beberapa substring, tidak ada yang cocok",
			s:       "hello world",
			substrs: []string{"foo", "bar", "baz"},
			want:    false,
			mock:    false,
		},
		{
			name:    "test dengan mock strings.Contains selalu false",
			s:       "hello world",
			substrs: []string{"world"},
			want:    false,
			mock:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var patches *gomonkey.Patches
			if tt.mock {
				patches = gomonkey.ApplyFuncReturn(strings.Contains, false)
				defer patches.Reset()
			}

			got := contains(tt.s, tt.substrs...)
			if got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
