package mail

import (
	"io"
	"net/mail"
)

type Mailer interface {
	Send(
		to, from mail.Address,
		subject string,
		body string,
		attachments map[string]io.Reader,
	) error
}
