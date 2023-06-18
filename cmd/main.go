package main

import (
	"log"

	todo "github.com/bb4ssttardio/RESTapi_todo-app"
)

func main() {
	srv := new(todo.Server)
	if err := srv.Run("8000"); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
