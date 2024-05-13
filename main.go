package main

import (
	"errors"
	"log"
	"os"

	"github.com/Davincible/goinsta/v3"
	"github.com/joho/godotenv"
)

func loadDotenv() (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", errors.New("error loading .env file")
	}
	username := os.Getenv("INSTAGRAM_USERNAME")
	password := os.Getenv("INSTAGRAM_PASSWORD")
	if username == "" || password == "" {
		return "", "", errors.New("username or password not found in .env file")
	}
	return username, password, nil
}

func main() {

	insta, err := goinsta.Import("./.goinsta")
	if err != nil {
		log.Println("Error loading .goinsta file, trying to login from dotenv")
		username, password, err := loadDotenv()
		if err != nil {
			log.Fatal(err)
		}
		insta = goinsta.New(username, password)
		if err := insta.Login(); err != nil {
			log.Fatal("Dotenv login unsuccessful; ", err)
		}
	} else {
		log.Println("Loaded .goinsta file")
	}

	defer insta.Export("./.goinsta")
}
