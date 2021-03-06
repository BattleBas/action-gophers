package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v29/github"
	"github.com/sethvargo/go-githubactions"
)

func main() {

	githubToken := githubactions.GetInput("GITHUB_TOKEN")
	if githubToken == "" {
		githubactions.Fatalf("missing input 'githubToken'")
		return
	}

	event, err := ioutil.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		log.Printf("Can't read events: " + err.Error())
		return
	}

	var pr github.PullRequestEvent
	json.Unmarshal(event, &pr)

	log.Printf("Number: %d Repo Name: %s, Owner: %s", *pr.PullRequest.Number, *pr.Repo.Name, *pr.Repo.Owner.Login)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// TODO make and add more Gopher gifs!
	url := "https://media.giphy.com/media/J0CZBmDGBStZQK4nfe/giphy.gif"

	message := "![GopherGif](" + url + ")"
	comment := &github.IssueComment{
		Body: &message,
	}
	_, _, err = client.Issues.CreateComment(ctx, *pr.Repo.Owner.Login, *pr.Repo.Name, *pr.PullRequest.Number, comment)
	if err != nil {
		log.Printf(err.Error())
		return
	}
}
