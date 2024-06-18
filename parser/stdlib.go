package parser

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
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

		respBody := responseWriter.body
		if respBody == nil {
			return
		}

		var resp any
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return
		}

		switch out := resp.(type) {
		case map[string]any:
			for k, v := range out {
				switch reflect.ValueOf(v).Kind() {
				case reflect.String:
					out[k] = ""
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
					reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
					out[k] = 0
				case reflect.Bool:
					out[k] = false
				case reflect.Interface:
					out[k] = map[string]string{}
				case reflect.Array:
				case reflect.Slice:
					out[k] = []int{}
				}
			}
		}

		b, err := json.Marshal(resp)
		if err != nil {
			return
		}

		log.Printf("%s %s\n", r.URL.String(), string(b))
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
