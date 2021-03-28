package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "send-me-a-secret"
	app.Usage = "distribute short secrets without a fuss"
	app.Authors = []*cli.Author{{
		Name:  "Robbie Ostrow",
		Email: "sendmeasecret@ostro.ws",
	}}
	app.Commands = []*cli.Command{
		{
			Name:    "initialize",
			Aliases: []string{"rotate"}, // TODO
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "skip-verify",
					Usage: "Skip verifying that the keys were persisted to disk and GitHub",
				},
				&cli.BoolFlag{
					Name:  "skip-github",
					Usage: "Don't upload public key to GitHub; instead, just print it out so it can be manually added. --skip-github implies --skip-verify",
				},
				&cli.StringFlag{
					Name:    "password",
					Usage:   "Password you'd like to use to encrypt your private key",
					EnvVars: []string{"SEND_ME_A_SECRET_PASSWORD"},
				},
				&cli.StringFlag{
					Name:    "github-token",
					Usage:   "GitHub token with write:public_key access. Ignored if --skip-github is specified.",
					EnvVars: []string{"GITHUB_TOKEN"},
				},
			},
			Usage: "Generate a new private/public key pair, saving the private key to disk and uploading the public key to GitHub",
			Action: func(c *cli.Context) error {
				return initialize(c.Bool("skip-verify"), c.Bool("skip-github"), c.String("password"), c.String("github-token"))
			},
		},
		{
			Name:      "decrypt",
			Aliases:   []string{"d"},
			Usage:     "Decrypt a string encrypted by `send-me-a-secret encrypt`. The string can be passed as an argument or via standard in.",
			ArgsUsage: "[string to decrypt]",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "password", Usage: "Password used to encrypt your private key", EnvVars: []string{"SEND_ME_A_SECRET_PASSWORD"}},
			},
			Action: func(c *cli.Context) error {
				return decrypt(c.String("password"), c.Args().First())
			},
		},
		{
			Name:      "encrypt",
			Aliases:   []string{"e"},
			Usage:     "Encrypt a short string so it can be decrypted by `send-me-a-secret decrypt`. The string can be passed as an argument or via standard in.",
			ArgsUsage: "[string to encrypt]",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Required: true, Usage: "The GitHub username of the person for whom this secret is intended"},
			},
			Action: func(c *cli.Context) error {
				return encrypt(c.String("user"), c.Args().First())
			},
		},
		{
			Name:  "validate",
			Usage: "Validate that someone can send you a secret",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Required: true, Usage: "Your GitHub username"},
				&cli.StringFlag{Name: "password", Usage: "Password used to encrypt your private key", EnvVars: []string{"SEND_ME_A_SECRET_PASSWORD"}},
			},
			Action: func(c *cli.Context) error {
				return validate(c.String("user"), c.String("password"))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
