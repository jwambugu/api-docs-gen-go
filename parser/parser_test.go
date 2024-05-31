package parser_test

import (
	"docs-gen/parser"
	"testing"
)

func TestParse(t *testing.T) {

	tests := []struct {
		filename string
		output   parser.Output
	}{
		{"test-docs", parser.OutputJSON},
		{"test-docs", parser.OutputYAML},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.filename, func(t *testing.T) {
			docsParser := parser.New(
				parser.WithPattern("*.go"),
				parser.WithOutput(tt.output),
				parser.WithFilename(tt.filename),
			)

			endpoints, err := docsParser.Parse()
			if err != nil {
				t.Fatalf("parse want no error, got %v", err)
			}

			if len(endpoints) != 2 {
				t.Fatalf("want 2 endpoints, got %v", len(endpoints))
			}

			if err = docsParser.ToFile(); err != nil {
				t.Fatalf("parser to file want no error, got %v", err)
			}
		})
	}
}
