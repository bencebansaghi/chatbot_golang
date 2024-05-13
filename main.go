package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Davincible/goinsta/v3"
	"github.com/joho/godotenv"
)

const (
	configPath = "./.goinsta"
)

func loadEnvVar(varName string) (string, error) {
	variable, ok := os.LookupEnv(varName)
	if !ok {
		return "", errors.New(varName + " not found in .env file")
	}
	if variable == "" {
		return "", errors.New(varName + " is empty")
	}
	return variable, nil
}

func loadEnvInstaUserPass() (string, string, error) {
	username, err := loadEnvVar("INSTAGRAM_USERNAME")
	if err != nil {
		return "", "", err
	}
	password, err := loadEnvVar("INSTAGRAM_PASSWORD")
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}

func getPosts(insta *goinsta.Instagram) /*[]goinsta.Item*/ {
	profilesStr, err := loadEnvVar("INSTA_PROFILES")
	if err != nil {
		log.Fatal(err)
	}
	profiles := strings.Split(profilesStr, ",")
	getPost(insta, profiles[len(profiles)-1])
}

func getPost(insta *goinsta.Instagram, profileStr string) /*goinsta.Item, */ {
	profile, err := insta.VisitProfile(profileStr)
	if err != nil {
		log.Fatal(err)
	}

	user := profile.User
	fmt.Printf(
		"%s has %d followers, %d posts, and %d IGTV vids\n",
		profileStr, user.FollowerCount, user.MediaCount, user.IGTVCount,
	)

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	insta, err := goinsta.Import(configPath)
	if err != nil {
		log.Println("Error loading .goinsta file, trying to login from dotenv")
		username, password, err := loadEnvInstaUserPass()
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
	err = insta.OpenApp()
	if err != nil {
		log.Fatal("Error opening app: ", err)
	}
	tl := insta.Timeline
	for _, item := range tl.Items {
		fmt.Println(item.Caption.Text)
	}

	// this dont work
	getPosts(insta)

	defer insta.Export(configPath)
}
