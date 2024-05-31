package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// _docsPrefix is the prefix used in comments to denote the start of API documentation.
const _docsPrefix = "@docs"

// _pathPrefix is the prefix used in comments to specify the API endpoint path.
const _pathPrefix = "@path"

// _methodPrefix is the prefix used in comments to specify the HTTP methods supported by the endpoint.
const _methodPrefix = "@method"

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

// _outputFilename is the default filename used to save the generated API documentation.
const _outputFilename = "docs.gen"

var htmlTmpl = ``

// Parameter represents an individual parameter for an API endpoint.
type Parameter struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

// Endpoint represents an API endpoint with its associated metadata.
// It includes the description, handler name, HTTP method, path, response type, and parameters.
type Endpoint struct {
	Description string      `json:"description" yaml:"description"`
	Handler     string      `json:"handler" yaml:"handler"`
	Method      string      `json:"method" yaml:"method"`
	Path        string      `json:"path" yaml:"path"`
	Response    string      `json:"response" yaml:"response"`
	Parameters  []Parameter `json:"parameters" yaml:"parameters"`
}

type Parser struct {
	endpoints []Endpoint
	filename  string
	files     []string
	output    Output
	pattern   string
}

// WithFilename is a functional option for configuring the filename used to save the generated API documentation.
func WithFilename(filename string) func(*Parser) {
	return func(p *Parser) {
		if p.filename != "" {
			p.filename = filename
		}
	}
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
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("parser: walk - %v", err)
		}

		if info.IsDir() {
			return nil
		}

		matches, err := filepath.Match(p.pattern, filepath.Base(path))
		if err != nil {
			return fmt.Errorf("parser: match - %v", err)
		}

		if matches {
			p.files = append(p.files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(p.files) == 0 {
		return nil, fmt.Errorf("parser: no files to parse")
	}

	for _, file := range p.files {
		fset := token.NewFileSet()

		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parser: parse file - %v", err)
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
						endpoint.Description = strings.TrimSpace(comment[len(_docsPrefix):])
					}

					if strings.HasPrefix(comment, _pathPrefix) {
						endpoint.Path = strings.TrimSpace(strings.TrimPrefix(comment, _pathPrefix))
					}

					if strings.HasPrefix(comment, _methodPrefix) {
						endpoint.Method = strings.TrimSpace(strings.TrimPrefix(comment, _methodPrefix))
					}

					if strings.HasPrefix(comment, _responsePrefix) {
						endpoint.Response = strings.TrimSpace(strings.TrimPrefix(comment, _responsePrefix))
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

				p.endpoints = append(p.endpoints, endpoint)
			}
			return true
		})

	}

	return p.endpoints, nil
}

// ToFile writes the generated API documentation to the specified file using the preferred Output format.
func (p *Parser) ToFile() error {
	if len(p.endpoints) == 0 {
		if _, err := p.Parse(); err != nil {
			return err
		}
	}

	var (
		outputBytes []byte
		err         error
	)

	switch p.output {
	case OutputJSON:
		outputBytes, err = json.Marshal(p.endpoints)
		if err != nil {
			return fmt.Errorf("parser: marshal - %v", err)
		}

		if ext := filepath.Ext(p.filename); ext != ".json" {
			p.filename += ".json"
		}

	case OutputYAML:
		outputBytes, err = yaml.Marshal(p.endpoints)
		if err != nil {
			return fmt.Errorf("parser: marshal - %v", err)
		}

		if ext := filepath.Ext(p.filename); ext != ".yml" && ext != ".yaml" {
			p.filename += ".yml"
		}
	}

	_, workingDir, _, _ := runtime.Caller(0)

	p.filename = filepath.Join(filepath.Dir(workingDir), "..", p.filename)

	if err = os.WriteFile(p.filename, outputBytes, 0644); err != nil {
		return fmt.Errorf("parser: write file - %v", err)
	}

	return nil
}

// New creates a new Parser instance based on the provided pattern.
// The pattern specifies the structure or format that the parser will use to interpret data.
// If the pattern is empty, it defaults to "*.go" to match Go source files.
func New(options ...func(*Parser)) *Parser {
	p := &Parser{
		filename: _outputFilename,
		output:   OutputJSON,
		pattern:  "*.go",
	}

	for _, option := range options {
		option(p)
	}

	return p
}
