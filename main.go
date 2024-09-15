package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var ctx = context.Background()

func make_valid_file(file string) string {
	to_replace := [...]string{":", "/", "\\"}
	for _, str := range to_replace {
		file = strings.ReplaceAll(file, str, "-")
	}
	file = strings.Trim(file, " ")
	return file
}

func main() {
	client, _ := reddit.NewReadonlyClient()
	posts, _, err := client.Subreddit.HotPosts(ctx, "wallpapers", &reddit.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	home, _ := os.UserHomeDir()
	wallpaperFolder := filepath.Join(home, "wallpapers")
	fmt.Println(wallpaperFolder)

	for _, post := range posts {
		extension := ""
		if strings.Contains(post.URL, ".png") {
			extension = ".png"
		}
		if strings.Contains(post.URL, ".jpg") {
			extension = ".jpg"
		}
		if strings.Contains(post.URL, ".jpeg") {
			extension = ".jpeg"
		}
		if extension == "" {
			continue
		}

		fmt.Println(post.Title)
		filename := make_valid_file(post.Title) + extension
		path := filepath.Join(wallpaperFolder, filename)
		fmt.Println(path)

		out, err := os.Create(path)
		if err != nil {
			log.Println(err)
			continue
		}

		resp, err := http.Get(post.URL)
		defer resp.Body.Close()
		if err != nil {
			log.Println("Error getting image")
			log.Println(err)
			continue
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Println("Can't copy file")
			log.Println(err)
		}

	}
}
