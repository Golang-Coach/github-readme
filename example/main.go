package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"

	"github.com/Golang-Coach/github-readme"
)

func main() {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "<!-- Your github access token -->"},
	)
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, src)
	githubClient := github.NewGithub(httpClient)

	content, err := githubClient.GetReadme(ctx, "parnurzeal", "gorequest")

	fmt.Println(content)
	fmt.Println(err)

}
