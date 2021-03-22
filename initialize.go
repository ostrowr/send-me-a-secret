// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The tokenauth command demonstrates using the oauth2.StaticTokenSource.
package main

import (
	"fmt"
	"syscall"

	"github.com/ostrowr/send-me-a-secret/internal/githubapi"
	"golang.org/x/term"
)

func main() {
	fmt.Print("GitHub Token: ")
	byteToken, _ := term.ReadPassword(int(syscall.Stdin))
	token := string(byteToken)
	client := githubapi.GetGithubClient(token)
	err := githubapi.CreateGist(client, "test", "content")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sdf")
}
