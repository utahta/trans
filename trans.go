package trans

import (
	"context"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type (
	// Client represents google translate api client
	Client struct {
		*translate.Client
	}
)

const (
	EnvTransAPIKey = "TRANS_API_KEY"
)

// New returns Client
func New(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	c, err := translate.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

// Translate translates input text
func (c *Client) Translate(ctx context.Context, input string, s string, t string, reverse bool) (string, error) {
	var (
		source language.Tag
		target language.Tag
		err    error
	)
	if s != "" {
		source, err = language.Parse(s)
		if err != nil {
			return "", err
		}
	}
	target, err = language.Parse(t)
	if err != nil {
		return "", err
	}

	var outputs []string
	for {
		res, err := c.Client.Translate(ctx, []string{input}, target, &translate.Options{Source: source, Format: translate.Text})
		if err != nil {
			return "", err
		}

		// trim \u200b (ZERO WIDTH SPACE)
		// http://unicode.org/cldr/utility/character.jsp?a=200B
		// https://www.fileformat.info/info/unicode/char/200B/index.htm
		outputs = append(outputs, strings.Replace(res[0].Text, "\u200b", "", -1))

		if reverse {
			reverse = false
			input = outputs[0]
			if source.IsRoot() {
				target = res[0].Source
			} else {
				target = source
			}
			source = language.Tag{}
			continue
		}
		break
	}

	return strings.Join(outputs, "\n"), nil
}
