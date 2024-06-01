package main

import "github.com/joho/godotenv"

func main() {

}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

}
