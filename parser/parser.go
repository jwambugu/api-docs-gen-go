package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// _docsPrefix is the prefix used in comments to denote the start of API documentation.
const _docsPrefix = "@docs"

// _pathPrefix is the prefix used in comments to specify the API endpoint path.
const _pathPrefix = "@path"

// _methodsPrefix is the prefix used in comments to specify the HTTP methods supported by the endpoint.
const _methodsPrefix = "@methods"

// _responsePrefix is the prefix used in comments to specify the response type of the endpoint.
const _responsePrefix = "@response"

// _parametersPrefix is the prefix used in comments to specify the parameters accepted by the endpoint.
const _parametersPrefix = "@parameters"

// Parameter represents an individual parameter for an API endpoint.
type Parameter struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

// Endpoint represents an API endpoint with its associated metadata.
// It includes the description, handler name, HTTP method, path, response type, and parameters.
type Endpoint struct {
	Description string      `json:"description"`
	Handler     string      `json:"handler"`
	Method      string      `json:"method"`
	Path        string      `json:"path"`
	Response    string      `json:"response"`
	Parameters  []Parameter `json:"parameters"`
}

type Parser struct {
	files   []string
	pattern string
}

// Parse analyzes the files matched by the parser's pattern and extracts API endpoint information.
// It returns a slice of Endpoint structs containing metadata about each API endpoint found.
func (p *Parser) Parse() ([]Endpoint, error) {
	var endpoints []Endpoint

	for _, file := range p.files {
		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse file: %v", err)
		}

		ast.Inspect(f, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Doc == nil {
					return false
				}

				commentsBlock := x.Doc.Text()
				if !strings.HasPrefix(commentsBlock, _docsPrefix) {
					return false
				}

				endpoint := Endpoint{
					Handler: x.Name.Name,
				}

				for _, comment := range strings.Split(commentsBlock, "\n") {
					if strings.HasPrefix(comment, _docsPrefix) {
						endpoint.Description = comment[len(_docsPrefix):]
					}

					if strings.HasPrefix(comment, _pathPrefix) {
						endpoint.Path = strings.TrimPrefix(comment, _pathPrefix)
					}

					if strings.HasPrefix(comment, _methodsPrefix) {
						endpoint.Method = strings.TrimPrefix(comment, _methodsPrefix)
					}

					if strings.HasPrefix(comment, _responsePrefix) {
						endpoint.Response = strings.TrimPrefix(comment, _responsePrefix)
					}

					if strings.HasPrefix(comment, _parametersPrefix) {
						fields := strings.Fields(strings.TrimPrefix(comment, _parametersPrefix))

						if len(fields) == 3 {
							endpoint.Parameters = append(endpoint.Parameters, Parameter{
								Name:     fields[0],
								Type:     fields[1],
								Required: fields[2] == "true",
							})
						}
					}
				}

				endpoints = append(endpoints, endpoint)
			}
			return true
		})

	}

	return endpoints, nil
}

// NewParser creates a new Parser instance based on the provided pattern.
// The pattern specifies the structure or format that the parser will use to interpret data.
// If the pattern is empty, it defaults to "*.go" to match Go source files.
func NewParser(pattern string) (*Parser, error) {
	if pattern == "" {
		pattern = "*.go"
	}

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob pattern: %v", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files to parse")
	}

	return &Parser{
		files:   files,
		pattern: pattern,
	}, nil
}
