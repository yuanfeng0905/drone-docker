package main

import (
	"strings"
	"testing"
)

func Test_fmtNow(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test now",
			want: "20210622",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmtNow(); !strings.HasPrefix(got, tt.want) {
				t.Errorf("fmtNow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genTag(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genTag(); got != tt.want {
				t.Errorf("genTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
