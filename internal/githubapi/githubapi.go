// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The tokenauth command demonstrates using the oauth2.StaticTokenSource.
package githubapi

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func CreateGist(githubClient *github.Client, name github.GistFilename, content string) error {
	description := "test"
	public := false
	files := make(map[github.GistFilename]github.GistFile)
	files[name] = github.GistFile{Content: &content}

	gist := &github.Gist{
		Description: &description,
		Public:      &public,
		Files:       files,
	}
	ctx := context.Background()
	_, _, err := githubClient.Gists.Create(ctx, gist)
	return err
}
