package main

import (
	"authorized/internal/server"
)

func init() {

}

func main() {
	srv := server.NewServer()
	panic(srv.Run("localhost:8000"))
}
