# Api Docs Generator

## Objective
Your assignment is to create an API documentation system using Go and any framework that will automatically generate documentation for your API endpoints.

### Help
Introspection can be useful for this challenge as it allows you to inspect the schema of your API at runtime.
With introspection, you can programmatically retrieve information about your API, such as the available endpoints, data types, and responses.
This can be useful for automatically generating documentation for your API and ensuring that it is always up-to-date.
Additionally, introspection can be used to implement features such as the ability to run requests against the API from the documentation,
as well as the ability to specify authentication credentials or API keys.

## Tasks

-   Implement assignment using:
    -   Language: Go
    -   Framework: any framework
-   Your system should automatically generate documentation for all of your API endpoints
-   The documentation should include:
    -   A list of available endpoints
    -   The HTTP methods supported by each endpoint
    -   The parameters each endpoint accepts, including their data type and whether they are required or optional
    -   The response format for each endpoint, including any error responses
-   Your system should be able to handle any changes to the API endpoints automatically
-   _Do not use pre-built API documentation tools like Swagger or API Blueprint_.

## TODO

- [x] Generate docs comments
- [x] Write doc contents to specified format, either `json` or `yml`
- [] Introspect at runtime
- [] Run requests, _maybe_
- [] Add CLI, _maybe_

## Usage

```go

package main

import (
	"docs-gen/parser"
	"log"
)

func main() {
	docsParser := parser.New(
		parser.WithPattern("*.go"),
		parser.WithOutput(parser.OutputJSON),
		parser.WithFilename("api-docs"),
	)

	// List all endpoints with docs
	endpoints, err := docsParser.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	for _, endpoint := range endpoints {
		log.Printf("[%s] %s - %s", endpoint.Method, endpoint.Path, endpoint.Description)
	}

	// 2024/05/31 13:56:40 [GET] /users - GetUsersHandler returns all users available.
	// 2024/05/31 13:56:40 [POST] /users - CreateUserHandler handles the creation of a new user.

	// docsParser.ToFile() // To export to file

}
```

### Sample Outputs

- json

```json

[
  {
    "description": "GetUsersHandler returns all users available.",
    "handler": "GetUsersHandler",
    "method": "GET",
    "path": "/users",
    "response": "GetUsersResponse",
    "parameters": null
  },
  {
    "description": "CreateUserHandler handles the creation of a new user.",
    "handler": "CreateUserHandler",
    "method": "POST",
    "path": "/users",
    "response": "CreateUserResponse",
    "parameters": [
      {
        "name": "name",
        "required": true,
        "type": "string"
      },
      {
        "name": "email",
        "required": true,
        "type": "string"
      }
    ]
  }
]
```

- yml

```yaml
- description: GetUsersHandler returns all users available.
  handler: GetUsersHandler
  method: GET
  path: /users
  response: GetUsersResponse
  parameters: [ ]
- description: CreateUserHandler handles the creation of a new user.
  handler: CreateUserHandler
  method: POST
  path: /users
  response: CreateUserResponse
  parameters:
    - name: name
      required: true
      type: string
    - name: email
      required: true
      type: string

```
