package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Davincible/goinsta/v3"
	"github.com/joho/godotenv"
)

const (
	configPath = "./.goinsta"
)

type Post struct {
	Caption   string
	ShortCode string
	Username  string
}

func loadEnvVar(varName string) (string, error) {
	variable, ok := os.LookupEnv(varName)
	if !ok {
		return "", errors.New(varName + " not found in .env file")
	} else if variable == "" {
		return "", errors.New(varName + " is empty")
	}
	return variable, nil
}

func loadEnvInstaUserPass() (string, string, error) {
	creds, err := loadEnvVar("INSTA_CREDS")
	if err != nil {
		return "", "", err
	}
	credsSplit := strings.Split(creds, ":")
	return credsSplit[0], credsSplit[1], nil
}

func getPosts(insta *goinsta.Instagram) []Post {
	profilesStr, err := loadEnvVar("INSTA_PROFILES")
	if err != nil {
		log.Fatal("posts: ", err)
	}

	profiles := strings.Split(profilesStr, ",")
	allPosts := []Post{}
	for _, profile := range profiles {
		posts := getPostsLastDay(insta, profile)
		if len(posts) != 0 {
			allPosts = append(allPosts, posts...)
		}
	}
	return allPosts
}

func getPostsLastDay(insta *goinsta.Instagram, profileStr string) []Post {
	profile, err := insta.VisitProfile(profileStr)
	if err != nil {
		log.Printf("[W]: Error while fetching %s profile: %s\n", profileStr, err)
		return []Post{}
	}

	lastDayPosts := []Post{}
	feed := profile.Feed
	if len(feed.Items) == 0 {
		log.Printf("[W]: No posts found for %s\n", profileStr)
		return []Post{}
	}

	for _, post := range feed.Items {
		// Check if the last post was taken within the last day
		if !time.Unix(post.TakenAt, 0).Before(time.Now().AddDate(0, 0, -1)) {
			lastDayPosts = append(lastDayPosts, Post{
				Caption:   post.Caption.Text,
				ShortCode: post.Code,
				Username:  profileStr,
			})
		}
	}
	// Only for testing purposes, should be removed later so it doesnt spam the logs
	log.Printf("[I]: Found %d posts for %s in the last day\n", len(lastDayPosts), profileStr)
	return lastDayPosts
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	insta, err := goinsta.Import(configPath, true) // skips initial sync as the profile is only used for the timeline, not for posting
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

	defer insta.Export(configPath)

	if insta.OpenApp() != nil {
		log.Fatal("Error opening app: ", err)
	}

	fmt.Println("Getting posts")
	fmt.Println(getPosts(insta))
	// Thats basically it after this this just goes directly to the python code
}
