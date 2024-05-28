package parser

import "net/http"

// @docs GetUsersHandler returns all users available.
// @path /users
// @method GET
// @response GetUsersResponse
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {}

// @docs CreateUserHandler handles the creation of a new user
// @path /users
// @method POST
// @parameters name string true
// @parameters email string true
// @response CreateUserResponse
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {}
