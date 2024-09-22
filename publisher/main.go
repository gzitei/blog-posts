package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

type Post struct {
	PostedAt string   `json:"posted-at"`
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	Image    string   `json:"image"`
}

type PostList struct {
	Posts []Post `json:"posts"`
}

func main() {
	dir := os.Args[1]

	metadataJson := path.Join(dir, "metadata.json")
	metadata, err := os.ReadFile(metadataJson)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	var posts PostList
	err = json.Unmarshal(metadata, &posts)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	done := []string{}
	for _, post := range posts.Posts {
		done = append(done, post.ID)
	}

	postsDir := path.Join(dir, "posts")
	directories, err := os.ReadDir(postsDir)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	for _, dir := range directories {
		dirName := path.Join("posts", dir.Name())
		if slices.Contains(done, dirName) {
			continue
		}

		files, err := os.ReadDir(dirName)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			return
		}

		post := Post{
			ID:   dirName,
			Tags: []string{},
		}

		for _, file := range files {
			fileName := file.Name()
			filePath := path.Join(dirName, fileName)
			if strings.HasPrefix(fileName, "#") {
				post.Tags = append(post.Tags, fileName)
			} else if strings.HasSuffix(fileName, "webp") {
				post.Image = fileName
			} else if fileName == "README.md" {
				file, err := os.Open(filePath)
				if err != nil {
					os.Stderr.WriteString(err.Error())
					return
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				if scanner.Scan() {
					firstLine := scanner.Text()
					title := strings.TrimLeft(firstLine, "# ")
					post.Title = strings.TrimSpace(title)
				}
			}
		}
		now := time.Now()
		strTime := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
		post.PostedAt = strTime
		posts.Posts = append(posts.Posts, post)
	}
	str, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	os.WriteFile(metadataJson, str, 0755)
}
