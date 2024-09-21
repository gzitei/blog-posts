package main

import (
	"encoding/json"
	"os"
	"path"
	"slices"
	"strings"
)

type Post struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	RawLink string   `json:"raw-link"`
	Tags    []string `json:"tags"`
	Image   string   `json:"image"`
}

type PostList struct {
	Posts []Post `json:"posts"`
}

func main() {
	dir := os.Args[1]

	const githubRawLink = "https://raw.githubusercontent.com/gzitei/blog-posts/refs/heads/main/"

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
		readme := path.Join(dirName, "README.md")
		post.RawLink = githubRawLink + readme

		for _, file := range files {
			fileName := file.Name()
			filePath := path.Join(dirName, fileName)
			if strings.HasPrefix(fileName, "#") {
				post.Tags = append(post.Tags, fileName)
			} else if strings.HasSuffix(fileName, "webp") {
				post.Image = githubRawLink + filePath
			} else if fileName == "README.md" {
				content, err := os.ReadFile(filePath)
				if err != nil {
					os.Stderr.WriteString(err.Error())
					return
				}
				title := strings.SplitN(string(content), "\n", 1)[0]
				title = strings.ReplaceAll(title, "#", "")
				post.Title = title
			}
		}
		posts.Posts = append(posts.Posts, post)
	}
	str, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	os.WriteFile(metadataJson, str, 0755)
}
