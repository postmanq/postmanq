package services

import (
	"context"
	"fmt"
	"github.com/postmanq/postmanq/pkg/plugins/smtpfx/smtp"
	"regexp"
	"strings"
)

var (
	rfc5322       = "(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
	rfc5322Regexp = regexp.MustCompile(rfc5322)
)

func NewFxEmailParser() smtp.EmailParser {
	return &emailParser{}
}

type emailParser struct {
}

func (p emailParser) Parse(ctx context.Context, email string) (*smtp.EmailAddress, error) {
	if !rfc5322Regexp.MatchString(email) {
		return nil, fmt.Errorf("%s has incorrect format", email)
	}

	i := strings.LastIndexByte(email, '@')
	return &smtp.EmailAddress{
		Address:   email,
		LocalPart: email[:i],
		Domain:    email[i+1:],
	}, nil
}
