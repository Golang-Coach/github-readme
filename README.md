# Github Readme
This library is used to get GitHub readme content in HTML as well as in JSON format

## How to use it
```go

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "<!-- Your github access token -->"},
	)
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, src)
	githubClient := github.NewGithub(httpClient)

	content, err := githubClient.GetReadme(ctx, "<!-- organization name -->", "<!-- Github repository name -->")

	fmt.Println(content)
	fmt.Println(err)

```
