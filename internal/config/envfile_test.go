package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveEnvFile_findsEnvTestInModuleRoot(t *testing.T) {
	path, err := ResolveEnvFile(DefaultEnvTestFile)
	if err != nil {
		t.Fatalf("ResolveEnvFile: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}

	root, err := ModuleRoot()
	if err != nil {
		t.Fatalf("ModuleRoot: %v", err)
	}

	if filepath.Base(path) != DefaultEnvTestFile {
		t.Fatalf("base = %s, want %s", filepath.Base(path), DefaultEnvTestFile)
	}
	if filepath.Dir(path) != root {
		t.Fatalf("dir = %s, root = %s", filepath.Dir(path), root)
	}
}

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
