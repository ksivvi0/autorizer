package main

import (
	"authorized/internal/server"
	"authorized/internal/services"
)

func init() {

}

func main() {
	l, err := services.NewLogger("./authorizer.log")
	if err != nil {
		panic(err)
	}
	a := services.NewAuth()

	servs := services.NewServices(l, a)
	srv := server.NewServer("localhost:8000", servs, true)
	panic(srv.Run("localhost:8000"))
}
