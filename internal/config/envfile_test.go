package config

import (
	"testing"
)

func TestEnvFileForAppEnv(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"dev", DefaultEnvDevFile},
		{"development", DefaultEnvDevFile},
		{"", DefaultEnvDevFile},
		{"prod", DefaultEnvProdFile},
		{"production", DefaultEnvProdFile},
		{"test", DefaultEnvTestFile},
		{"staging", ".env.staging"},
	}
	for _, tt := range tests {
		if got := EnvFileForAppEnv(tt.in); got != tt.want {
			t.Fatalf("EnvFileForAppEnv(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
