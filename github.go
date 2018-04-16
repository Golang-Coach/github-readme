package github

import (
	"context"
	"fmt"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Github struct {
	client *http.Client
}

// NewGithub initialized Github
func NewGithub(client *http.Client) *Github {
	return &Github{
		client: client,
	}
}

func withContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

// GetReadme returns readme content in HTML as well as in JSON format
func (g *Github) GetReadme(ctx context.Context, owner string, repo string) (string, error) {
	u := fmt.Sprintf("repos/%v/%v/readme", owner, repo)

	req, err := http.NewRequest("GET", "https://api.github.com/"+u, nil) // TODO -- need to refactored
	req = withContext(ctx, req)
	req.Header.Add("Accept", "application/vnd.github.v3.html+json")
	resp, err := g.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context'g error is probably more useful.
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return "", e
			}
		}

		return "", err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	if c := resp.StatusCode; 200 <= c && c <= 299 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(data[:]), nil
	}
	return "", errors.New(resp.Status)
}
