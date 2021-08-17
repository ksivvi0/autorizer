package main

import (
	"authorizer/internal/server"
	"authorizer/internal/services"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	err := godotenv.Load(".env") //full path
	if err != nil {
		panic(err)
	}
}

func main() {

	l, err := services.NewLoggerInstance(os.Getenv("LOG_PATH"))
	if err != nil {
		panic(err)
	}
	a := services.NewAuthInstance()
	s, err := services.NewStoreInstance(os.Getenv("MONGO_URI"))
	if err != nil {
		writeLogNPanic(l, err)
	}
	if err = s.Ping(); err != nil {
		writeLogNPanic(l, err)
	}

	svcs := services.NewServices(l, a, s)
	srv, err := server.NewServerInstance(os.Getenv("SERVER_ADDR"), svcs, true)
	if err != nil {
		writeLogNPanic(l, err)
	}
	panic(srv.Run())
}

func writeLogNPanic(l *services.Logger, err error) {
	l.WriteError(err.Error())
	panic(err)
}
