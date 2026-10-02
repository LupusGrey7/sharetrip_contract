package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	yaml3 "go.yaml.in/yaml/v3"
)

func TestOrderPathsPreservesStaticRoutesBeforeParameters(t *testing.T) {
	input := []byte(`
openapi: 3.0.3
paths:
  /contracts/{contractId}: {}
  /contracts/active: {}
  /contracts/openapi: {}
`)
	order := []string{
		"/contracts/active",
		"/contracts/openapi",
		"/contracts/{contractId}",
	}

	output, err := orderPaths(input, order)
	if err != nil {
		t.Fatal(err)
	}

	var document yaml3.Node
	if err := yaml3.Unmarshal(output, &document); err != nil {
		t.Fatal(err)
	}

	root := document.Content[0]
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value != "paths" {
			continue
		}

		paths := root.Content[i+1]
		got := make([]string, 0, len(paths.Content)/2)
		for j := 0; j < len(paths.Content); j += 2 {
			got = append(got, paths.Content[j].Value)
		}
		if !reflect.DeepEqual(got, order) {
			t.Fatalf("path order = %v, want %v", got, order)
		}
		return
	}

	t.Fatal("paths mapping not found")
}

func TestOpenAPIPathOrderReadsSourceOrder(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "contract.yaml")
	content := []byte("openapi: 3.0.3\npaths:\n  /contracts/active: {}\n  /contracts/{contractId}: {}\n")
	if err := os.WriteFile(filename, content, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := openAPIPathOrder(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/contracts/active", "/contracts/{contractId}"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("path order = %v, want %v", got, want)
	}
}
