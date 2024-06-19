package parser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// @docs GetUsersHandler returns all users available.
// @path /users
// @method GET
// @response GetUsersResponse
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(http.StatusText(http.StatusCreated)))
}

func TestParse(t *testing.T) {
	for _, output := range []Output{OutputJSON, OutputYAML} {
		t.Run(string(output), func(t *testing.T) {
			f, err := os.CreateTemp("", "docs")
			require.NoError(t, err)

			t.Cleanup(func() {
				err = os.Remove(f.Name())
				require.NoError(t, err)
			})

			docsParser := New(
				WithPattern("*.go"),
				WithOutput(output),
				WithFilename(f.Name()),
			)

			endpoints, err := docsParser.Parse()
			require.NoError(t, err)
			require.Len(t, endpoints, 2)

			err = docsParser.ToFile()
			require.NoError(t, err)

			contents, err := os.ReadFile(f.Name())
			require.NoError(t, err)
			require.NotNil(t, contents)
		})
	}
}

func TestParser_StdLib(t *testing.T) {
	file, err := os.CreateTemp("", "*.json")
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.Remove(file.Name())
		require.NoError(t, err)
	})

	var (
		docsParser = New(WithFilename(file.Name()))
		router     = http.NewServeMux()
	)

	router.Handle("GET /users", StdLib(GetUsersHandler, docsParser))

	var (
		req = httptest.NewRequest(http.MethodGet, "/users", nil)
		rr  = httptest.NewRecorder()
	)

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp User
	err = json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)

	err = docsParser.ToFile()
	require.NoError(t, err)

	contents, err := os.ReadFile(file.Name())
	require.NoError(t, err)

	var endpoints map[string]Endpoint
	err = json.Unmarshal(contents, &endpoints)
	require.NoError(t, err)

	wantResponse := map[string]any{"name": ""}
	endpoint := endpoints[getKey(http.MethodGet, "/users")]
	require.Equal(t, wantResponse, endpoint.Responses[0].Body)
	require.Equal(t, http.StatusOK, endpoint.Responses[0].StatusCode)
}
