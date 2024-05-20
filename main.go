package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
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
	username, err := loadEnvVar("INSTA_USERNAME")
	if err != nil {
		return "", "", err
	}
	password, err := loadEnvVar("INSTA_PASSWORD")
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}

func getPosts(insta *goinsta.Instagram, postsChan chan<- []Post) {
	profilesStr, err := loadEnvVar("INSTAGRAM_PAGES")
	if err != nil {
		log.Fatal("posts: ", err)
	}

	profiles := strings.Split(profilesStr, ",")
	var wg sync.WaitGroup
	for _, profile := range profiles {
		wg.Add(1)
		go func(profile string) {
			defer wg.Done()
			getPostsLastDay(insta, profile, postsChan)
		}(profile)
	}
	go func() {
		wg.Wait()
		close(postsChan)
	}()
}

func getPostsLastDay(insta *goinsta.Instagram, profileStr string, postsChan chan<- []Post) {
	profile, err := insta.VisitProfile(profileStr)
	if err != nil {
		log.Printf("[W]: Error while fetching %s profile: %s\n", profileStr, err)
		return
	}

	lastDayPosts := []Post{}
	feed := profile.Feed
	if len(feed.Items) == 0 {
		log.Printf("[W]: %s hasn't posted yet.\n", profileStr)
		return
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
	postsChan <- lastDayPosts
}

func main() {
	log.SetOutput(os.Stderr)
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

	postsChan := make(chan []Post)
	go getPosts(insta, postsChan)

	for posts := range postsChan {
		for _, post := range posts {
			jsonPost, err := json.Marshal(post)
			if err != nil {
				log.Printf("[W]: Error marshalling post: %s\n", err)
				continue
			}
			os.Stdout.Write(jsonPost)
			os.Stdout.WriteString("\n")
		}
	}
}
