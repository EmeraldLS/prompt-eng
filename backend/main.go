package main

import (
	"github.com/emeraldls/fyp/internal/rest"
	"github.com/joho/godotenv"
)

func main() {
	srv := rest.NewServer(rest.Config{
		ListAddr: ":2323",
	})

	srv.SetupRouter()
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

}
