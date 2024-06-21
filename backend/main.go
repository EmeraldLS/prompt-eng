package main

import (
	"github.com/emeraldls/fyp/internal/rest"
	"github.com/joho/godotenv"
)

func main() {
	srv := rest.NewServer(rest.Config{
		ListAddr:           ":2323",
		EntityExtractorURL: "http://localhost:2525/extract-entities",
		ClientOrigin:       "http://localhost:5173",
	})

	srv.SetupRouter()
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

}
