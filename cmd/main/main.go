package main

import (
	"authorized/internal/server"
	"authorized/internal/services"
)

func init() {

}

func main() {
	l, err := services.NewLogger("./log.")
	if err != nil {
		panic(err)
	}
	servs := services.NewServices(l)
	srv := server.NewServer(servs)
	panic(srv.Run("localhost:8000"))
}
