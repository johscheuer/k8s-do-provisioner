package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: t.AccessToken,
	}, nil
}

func createNewDOClient() *godo.Client {
	pat, err := ioutil.ReadFile("./.token")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	return godo.NewClient(oauth2.NewClient(
		oauth2.NoContext,
		&TokenSource{
			AccessToken: strings.TrimSpace(string(pat)),
		},
	))
}
