package main

import (
	"testing"
)

// TestFilteringNop ensures with no filtering nothing changes
func TestFilteringNop(t *testing.T) {

	in := []CodeBlock{
		CodeBlock{Name: "steve", Shell: "sh", Content: []string{"steve"}},
		CodeBlock{Name: "steve2", Shell: "sh2", Content: []string{"steve2"}},
	}

	sh := ""
	nm := ""
	shellArg = &sh
	nameArg = &nm

	out := filterBlocks(in)

	if len(out) != len(in) {
		t.Fatalf("unexepected filtering")
	}
}

// TestFilteringName ensures with Name matching we work
func TestFilteringName(t *testing.T) {

	in := []CodeBlock{
		CodeBlock{Name: "steve", Shell: "sh", Content: []string{"steve"}},
		CodeBlock{Name: "steve2", Shell: "sh2", Content: []string{"steve2"}},
	}

	sh := ""
	nm := "steve"
	shellArg = &sh
	nameArg = &nm

	out := filterBlocks(in)

	if len(out) != 1 {
		t.Fatalf("unexepected filtering")
	}
	if out[0].Content[0] != "steve" {
		t.Fatalf("wrong filtering")
	}
	if out[0].Shell != "sh" {
		t.Fatalf("wrong filtering")
	}
}

// TestFilteringShell ensures with Shell matching we work
func TestFilteringShell(t *testing.T) {

	in := []CodeBlock{
		CodeBlock{Name: "steve", Shell: "/bin/bash", Content: []string{"steve"}},
		CodeBlock{Name: "steve2", Shell: "/bin/sh", Content: []string{"steve2"}},
	}

	sh := "bash"
	nm := ""
	shellArg = &sh
	nameArg = &nm

	out := filterBlocks(in)

	if len(out) != 1 {
		t.Fatalf("unexepected filtering")
	}
	if out[0].Content[0] != "steve" {
		t.Fatalf("wrong filtering")
	}
	if out[0].Shell != "/bin/bash" {
		t.Fatalf("wrong filtering")
	}
}

// TestParser ensures we can parse our README.md file
func TestParser(t *testing.T) {

	x, err := parseBlocks("README.md")
	if err != nil {
		t.Fatalf("unexpected error parsing README.md")
	}

	expected := 1
	if len(x) != expected {
		t.Fatalf("got %d blocks, expected %d", len(x), expected)
	}

	_, err = parseBlocks("README.md.missing")
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
}

// TestExists is a trivial test of fileExists
func TestExists(t *testing.T) {

	a := fileExists("README.md")
	b := fileExists("README.md.missing")

	if !a {
		t.Fatalf("expected file to exist, but it doesnt")
	}

	if b {
		t.Fatalf("expected file to not exist, but it is present")
	}
}
