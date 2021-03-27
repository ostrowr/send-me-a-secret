// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The tokenauth command demonstrates using the oauth2.StaticTokenSource.
package githubapi

import (
	"context"
	"crypto/rsa"
	"errors"

	"github.com/google/go-github/v33/github"
	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
	"golang.org/x/oauth2"
)

// This is just a random string used as the public key name
// so it doesn't accidentally collide with other public keys that
// the user might have. Don't change this string!
var PublicKeyName = "send-me-a-secret-jwpyl"

func GetGithubClient(token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func UploadPublicKeyToGithub(githubClient *github.Client, publicKey string) error {
	ctx := context.Background()
	readonly := true
	key := github.Key{
		Title:    &PublicKeyName,
		Key:      &publicKey,
		ReadOnly: &readonly, // readonly doesn't seem to work
	}
	_, _, err := githubClient.Users.CreateKey(ctx, &key) // todo ensure that the key can't be created twice
	if err != nil {
		return err
	}
	return nil
}

var ErrNoValidKeys = errors.New("no valid keys found for user")
var ErrTooManyValidKeys = errors.New("multiple valid keys found for user; can't decide")

// GetPublicKeyFromGithub fetches an rsa key from `username` on github with the correct keyLength.
// (keyLength is hackily used to identify keys generated for send-me-a-secret, since the title is not publicly available.)
// If there are no valid keys or if there is more than one maching key, this returns an error.
func GetPublicKeyFromGithubUnauthenticated(githubClient *github.Client, username string, isValidKey func(*rsa.PublicKey) bool) (*rsa.PublicKey, error) {
	ctx := context.Background()
	options := github.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	allPublicKeys := make([]*github.Key, 0, 10)
	for {
		// TODO deal with rate limiting; probably in GetGithubClient
		keys, response, err := githubClient.Users.ListKeys(ctx, username, &options)
		if err != nil {
			return nil, err
		}
		allPublicKeys = append(allPublicKeys, keys...)
		if response.NextPage <= options.Page {
			break
		}
		options.Page = response.NextPage
	}

	validKeys := make([]*rsa.PublicKey, 0, 1)

	for _, key := range allPublicKeys {
		parsed, err := rsahelpers.SSHPubKeyToRSAPubKey([]byte(*key.Key))
		if err != nil {
			continue // it's ok if they have some public keys that we can't parse
		}
		if isValidKey(parsed) {
			validKeys = append(validKeys, parsed)
		}
	}

	if len(validKeys) == 0 {
		return nil, ErrNoValidKeys
	}
	if len(validKeys) > 1 {
		return nil, ErrTooManyValidKeys
	}
	return validKeys[0], nil
}

func GetCurrentUsername(githubClient *github.Client) (string, error) {
	ctx := context.Background()
	user, _, err := githubClient.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}
	return *user.Login, nil
}
