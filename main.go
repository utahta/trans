package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type Options struct {
	// ShowHelp shows help
	ShowHelp bool

	// Source is the language of the input strings.
	Source string

	// Target is the language of the output strings.
	Target string

	// CredentialsFile is the service account JSON credentials file
	// refs https://cloud.google.com/iam/docs/creating-managing-service-accounts
	CredentialsFile string
}

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
	flag.StringVar(&opts.CredentialsFile, "c", "", "Sets the service account JSON credentials file")
	flag.Parse()
	if opts.ShowHelp {
		flag.PrintDefaults()
		return nil
	}
	if opts.Target == "" {
		fmt.Println("The -t option is required")
		return nil
	}
	inputs := flag.Args()

	ctx := context.Background()
	var clientOpts []option.ClientOption
	if opts.CredentialsFile != "" {
		clientOpts = append(clientOpts, option.WithCredentialsFile(opts.CredentialsFile))
	}

	c, err := translate.NewClient(ctx, clientOpts...)
	if err != nil {
		return err
	}

	var (
		source language.Tag
		target language.Tag
	)
	if opts.Source != "" {
		source, err = language.Parse(opts.Source)
		if err != nil {
			return err
		}
	}
	target, err = language.Parse(opts.Target)
	if err != nil {
		return err
	}

	trans, err := c.Translate(ctx, inputs, target, &translate.Options{Source: source, Format: translate.Text})
	if err != nil {
		return err
	}

	for _, t := range trans {
		fmt.Println(t.Text)
	}
	return nil
}
