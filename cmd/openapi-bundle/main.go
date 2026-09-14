package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oasdiff/yaml"
	yaml3 "go.yaml.in/yaml/v3"
)

func main() {
	if len(os.Args) != 3 {
		_, err := fmt.Fprintln(os.Stderr, "usage: openapi-bundle <input> <output>")
		if err != nil {
			return
		}
		os.Exit(2)
	}

	ctx := context.Background()
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	pathOrder, err := openAPIPathOrder(os.Args[1])
	if err != nil {
		exitError("read OpenAPI path order", err)
	}
	doc, err := loader.LoadFromFile(os.Args[1])
	if err != nil {
		exitError("load OpenAPI", err)
	}

	doc.InternalizeRefs(ctx, preserveComponentName)
	if err := doc.Validate(ctx); err != nil {
		exitError("validate OpenAPI", err)
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		exitError("marshal OpenAPI", err)
	}
	data, err = yaml.JSONToYAML(data)
	if err != nil {
		exitError("convert OpenAPI to YAML", err)
	}
	data, err = orderPaths(data, pathOrder)
	if err != nil {
		exitError("preserve OpenAPI path order", err)
	}

	if err := os.MkdirAll(filepath.Dir(os.Args[2]), 0o755); err != nil {
		exitError("create output directory", err)
	}
	if err := os.WriteFile(os.Args[2], append(data, '\n'), 0o644); err != nil {
		exitError("write bundled OpenAPI", err)
	}
}

func openAPIPathOrder(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var document yaml3.Node
	if err := yaml3.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	root := document.Content[0]
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value != "paths" {
			continue
		}

		paths := root.Content[i+1]
		order := make([]string, 0, len(paths.Content)/2)
		for j := 0; j < len(paths.Content); j += 2 {
			order = append(order, paths.Content[j].Value)
		}
		return order, nil
	}

	return nil, fmt.Errorf("paths mapping not found")
}

func orderPaths(data []byte, pathOrder []string) ([]byte, error) {
	var document yaml3.Node
	if err := yaml3.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	root := document.Content[0]
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value != "paths" {
			continue
		}
		reorderMapping(root.Content[i+1], pathOrder)
		break
	}

	return yaml3.Marshal(&document)
}

func reorderMapping(mapping *yaml3.Node, order []string) {
	pairs := make(map[string][2]*yaml3.Node, len(mapping.Content)/2)
	for i := 0; i < len(mapping.Content); i += 2 {
		pairs[mapping.Content[i].Value] = [2]*yaml3.Node{mapping.Content[i], mapping.Content[i+1]}
	}

	mapping.Content = mapping.Content[:0]
	for _, key := range order {
		pair, ok := pairs[key]
		if !ok {
			continue
		}
		mapping.Content = append(mapping.Content, pair[0], pair[1])
		delete(pairs, key)
	}
	for _, pair := range pairs {
		mapping.Content = append(mapping.Content, pair[0], pair[1])
	}
}

func preserveComponentName(doc *openapi3.T, ref openapi3.ComponentRef) string {
	refPath := ref.RefPath()
	if refPath != nil {
		parts := strings.Split(strings.Trim(refPath.Fragment, "/"), "/")
		if len(parts) >= 3 && parts[len(parts)-3] == "components" {
			return parts[len(parts)-1]
		}
	}
	return openapi3.DefaultRefNameResolver(doc, ref)
}

func exitError(action string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", action, err)
	os.Exit(1)
}
