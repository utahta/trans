package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/utahta/trans"
	"google.golang.org/api/option"
)

type Options struct {
	// ShowHelp shows help
	ShowHelp bool

	// Source is the language of the input strings.
	Source string

	// Target is the language of the output strings.
	Target string

	// APIKey is the key of authenticate to the Translation API.
	// refs https://cloud.google.com/translate/docs/auth#using_an_api_key
	APIKey string

	// CredentialsFile is the service account JSON credentials file
	// refs https://cloud.google.com/iam/docs/creating-managing-service-accounts
	CredentialsFile string
}

const (
	envTransAPIKey = "TRANS_API_KEY"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	opts := &Options{}
	flag.BoolVar(&opts.ShowHelp, "h", false, "Show help")
	flag.StringVar(&opts.Source, "s", "", "Sets the source language")
	flag.StringVar(&opts.Target, "t", "", "Sets the target language")
	flag.StringVar(&opts.APIKey, "key", "", "Sets the api key")
	flag.StringVar(&opts.CredentialsFile, "c", "", "Sets the service account JSON credentials file")
	flag.Parse()
	if opts.ShowHelp {
		fmt.Println("Usage: trans -t ja TEXT")
		flag.PrintDefaults()
		return nil
	}
	if opts.Target == "" {
		fmt.Println("The -t option is required")
		return nil
	}
	input := strings.Join(flag.Args(), " ")

	ctx := context.Background()
	var clientOpts []option.ClientOption
	if opts.APIKey != "" {
		clientOpts = append(clientOpts, option.WithAPIKey(opts.APIKey))
	} else if os.Getenv(envTransAPIKey) != "" {
		clientOpts = append(clientOpts, option.WithAPIKey(os.Getenv(envTransAPIKey)))
	}

	if opts.CredentialsFile != "" {
		clientOpts = append(clientOpts, option.WithCredentialsFile(opts.CredentialsFile))
	}

	c, err := trans.New(ctx, clientOpts...)
	if err != nil {
		return err
	}

	text, err := c.Translate(ctx, input, opts.Source, opts.Target)
	if err != nil {
		return err
	}

	fmt.Println(text)
	return nil
}
