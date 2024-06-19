package main

import (
	"docs-gen/parser"
	"encoding/json"
	"log"
	"reflect"
)

func main() {

	parser.Srv()
	return
	type User struct {
		Email string `json:"email,omitempty"`
		ID    int64  `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		IDs   any
	}

	user := User{
		Name:  "Joram",
		ID:    1,
		Email: "j@user.com",
		IDs:   []int{1, 2, 3},
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(b))

	var resp any
	if err = json.Unmarshal(b, &resp); err != nil {
		log.Fatalln(err)
	}

	log.Printf("%T\n", resp)

	keys := reflect.ValueOf(resp).MapKeys()
	log.Println(keys)

	switch out := resp.(type) {
	case map[string]any:
		for k, v := range out {
			val := reflect.ValueOf(v)
			switch val.Kind() {
			case reflect.String:
				out[k] = ""
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				out[k] = 0
			case reflect.Bool:
				out[k] = false
			case reflect.Interface:
				out[k] = nil
			case reflect.Array:
			case reflect.Slice:
				out[k] = []int{}
			}
		}
	}

	b, _ = json.Marshal(resp)
	log.Println(string(b))
	// parser.Srv()
}
