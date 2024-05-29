package parser

import "testing"

func TestParse(t *testing.T) {
	parser := New(WithPattern("*.go"), WithOutput(OutputHTML))

	endpoints, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse want no error, got %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("want 2 endpoints, got %v", len(endpoints))
	}
}
