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

// Output represents the format in which the API documentation can be generated.
type Output string

const (
	// OutputJSON specifies that the API documentation should be generated in JSON format.
	OutputJSON Output = "json"

	// OutputYAML specifies that the API documentation should be generated in YAML format.
	OutputYAML Output = "yaml"

	// OutputHTML specifies that the API documentation should be generated in HTML format.
	OutputHTML Output = "html"
)

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
	output  Output
}

// WithPattern is a functional option for configuring the pattern used by a Parser.
// If the provided pattern is empty, it defaults to "*.go".
func WithPattern(pattern string) func(*Parser) {
	return func(p *Parser) {
		if pattern != "" {
			p.pattern = pattern
		}
	}
}

// WithOutput is a functional option for configuring the output format used by a Parser.
// If the provided output format is empty, it defaults to OutputJSON.
func WithOutput(output Output) func(*Parser) {
	return func(p *Parser) {
		if output != "" {
			p.output = output
		}
	}
}

// Parse analyzes the files matched by the parser's pattern and extracts API endpoint information.
// It returns a slice of Endpoint structs containing metadata about each API endpoint found.
func (p *Parser) Parse() ([]Endpoint, error) {
	var (
		endpoints []Endpoint
		err       error
	)

	p.files, err = filepath.Glob(p.pattern)
	if err != nil {
		return nil, fmt.Errorf("parser: glob pattern: %w", err)
	}

	if len(p.files) == 0 {
		return nil, fmt.Errorf("parser: no files to parse")
	}

	for _, file := range p.files {
		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parser: parse file: %v", err)
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

// New creates a new Parser instance based on the provided pattern.
// The pattern specifies the structure or format that the parser will use to interpret data.
// If the pattern is empty, it defaults to "*.go" to match Go source files.
func New(options ...func(*Parser)) *Parser {
	p := &Parser{
		pattern: "*.go",
		output:  OutputJSON,
	}

	for _, option := range options {
		option(p)
	}

	return p
}
