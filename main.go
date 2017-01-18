package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("app", Version)
	c.Args = os.Args[1:]

	// Grab the email from Environment variable
	if os.Getenv("UKC_EMAIL") != "" {
		c.Args = append(c.Args, "-email="+os.Getenv("UKC_EMAIL"))
	}

	// Grab the password from Environment variable
	if os.Getenv("UKC_PASSWORD") != "" {
		c.Args = append(c.Args, "-password="+os.Getenv("UKC_PASSWORD"))
	}

	// Override the default CLI version stuff
	for _, arg := range c.Args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(c.Args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], c.Args)
			c.Args = newArgs
			break
		}
	}

	c.Commands = Commands

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
