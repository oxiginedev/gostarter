package mail

import (
	"fmt"
	"io"
	"net/mail"
	"net/smtp"

	"github.com/domodwyer/mailyak/v3"
	"github.com/oxiginedev/gostarter/config"
	"github.com/oxiginedev/gostarter/util"
)

type SmtpMailer struct {
	client *mailyak.MailYak
}

func NewSmtpMailer(opts *config.MailConfiguration) (*SmtpMailer, error) {
	smtpAuth := smtp.PlainAuth("", opts.Username, opts.Password, opts.Host)

	yak := mailyak.New(fmt.Sprintf("%s:%d", opts.Host, opts.Port), smtpAuth)
	if opts.SSL {
		var tlsErr error
		yak, tlsErr = mailyak.NewWithTLS(fmt.Sprintf("%s:%d", opts.Host, opts.Port), smtpAuth, nil)
		if tlsErr != nil {
			return nil, tlsErr
		}
	}

	return &SmtpMailer{client: yak}, nil
}

func (s *SmtpMailer) Send(to, from mail.Address, subject, body string, attachments map[string]io.Reader) error {
	if !util.IsStringEmpty(from.Name) {
		s.client.FromName(from.Name)
	}

	s.client.From(from.Address)
	s.client.To(to.Address)
	s.client.Subject(subject)
	s.client.HTML().Set(body)

	for name, attachment := range attachments {
		s.client.Attach(name, attachment)
	}

	err := s.client.Send()
	if err != nil {
		return err
	}

	return nil
}
