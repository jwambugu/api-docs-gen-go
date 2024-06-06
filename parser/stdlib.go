package parser

import (
	"encoding/json"
	"log"
	"net/http"
)

type Middleware func(handler http.Handler) http.Handler

type stdResponseWriter struct {
	body       []byte
	statusCode int
	w          http.ResponseWriter
}

// Header returns the header map that will be sent by ResponseWriter.WriteHeader.
func (std *stdResponseWriter) Header() http.Header {
	return std.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (std *stdResponseWriter) Write(bytes []byte) (int, error) {
	std.body = bytes
	return std.w.Write(bytes)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (std *stdResponseWriter) WriteHeader(statusCode int) {
	std.statusCode = statusCode
	std.w.WriteHeader(statusCode)
}

func StdLib(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseWriter := &stdResponseWriter{w: w}
		next.ServeHTTP(responseWriter, r)

		if body := responseWriter.body; body != nil {
			var resp any
			if err := json.Unmarshal(body, &resp); err != nil {
				log.Printf("decode body: %v", err)
				return
			}

			log.Printf("%T", resp)
		}

	}
}

// @docs GetUsersHandler returns all users available.
// @path /users
// @method GET
// @response GetUsersResponse
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}

type User struct {
	Name string `json:"name"`
}

// @docs CreateUserHandler handles the creation of a new user.
// @path /users
// @method POST
// @parameters name string true
// @parameters email string true
// @response CreateUserResponse
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	resp := User{Name: "Jay"}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func Srv() {
	router := http.NewServeMux()
	//router.HandleFunc("GET /users", StdLib(GetUsersHandler))
	router.HandleFunc("GET /users", StdLib(CreateUserHandler))

	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalln(err)
	}
}
