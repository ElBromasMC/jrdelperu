package service

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/wneessen/go-mail"
)

var (
	ErrEmailNotConfigured = errors.New("servicio de email no configurado")
	ErrEmailNoRecipients  = errors.New("no se especificaron destinatarios")
)

// emailDialTimeout es el tiempo máximo para conectar y enviar un correo.
const emailDialTimeout = 15 * time.Second

// EmailAttachment representa un archivo adjunto a un correo.
type EmailAttachment struct {
	Filename    string
	Content     []byte
	ContentType string // ej: "application/pdf"
}

// Email representa un correo a enviar.
type Email struct {
	To          []string
	Subject     string
	HTMLBody    string
	TextBody    string // alternativa en texto plano (opcional)
	Attachments []EmailAttachment
}

// EmailSender es la interfaz de envío de correos, para poder intercambiar la
// implementación (SMTP, mocks en pruebas, etc.) fácilmente.
type EmailSender interface {
	Send(ctx context.Context, email Email) error
	IsConfigured() bool
	CompanyEmail() string
}

// EmailService envía correos mediante SMTP usando wneessen/go-mail.
type EmailService struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
	fromName  string
	toEmail   string
}

// NewEmailService crea una nueva instancia del servicio de correo.
func NewEmailService(host string, port int, username, password, fromEmail, fromName, toEmail string) *EmailService {
	return &EmailService{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
		fromName:  fromName,
		toEmail:   toEmail,
	}
}

// IsConfigured indica si el servicio tiene la configuración mínima para enviar.
func (s *EmailService) IsConfigured() bool {
	return s.host != "" && s.port != 0 && s.fromEmail != ""
}

// CompanyEmail devuelve el correo destinatario por defecto (SMTP_TO_EMAIL).
func (s *EmailService) CompanyEmail() string {
	return s.toEmail
}

// Send envía un correo. Devuelve ErrEmailNotConfigured si falta configuración.
func (s *EmailService) Send(ctx context.Context, email Email) error {
	if !s.IsConfigured() {
		return ErrEmailNotConfigured
	}
	if len(email.To) == 0 {
		return ErrEmailNoRecipients
	}

	msg := mail.NewMsg()
	if err := msg.FromFormat(s.fromName, s.fromEmail); err != nil {
		return err
	}
	if err := msg.To(email.To...); err != nil {
		return err
	}
	msg.Subject(email.Subject)

	// Cuerpo: si hay texto plano, se usa como principal y el HTML como alternativa.
	switch {
	case email.TextBody != "" && email.HTMLBody != "":
		msg.SetBodyString(mail.TypeTextPlain, email.TextBody)
		msg.AddAlternativeString(mail.TypeTextHTML, email.HTMLBody)
	case email.HTMLBody != "":
		msg.SetBodyString(mail.TypeTextHTML, email.HTMLBody)
	default:
		msg.SetBodyString(mail.TypeTextPlain, email.TextBody)
	}

	// Adjuntos (ej: la reclamación en PDF).
	for _, att := range email.Attachments {
		contentType := att.ContentType
		if contentType == "" {
			contentType = string(mail.TypeAppOctetStream)
		}
		if err := msg.AttachReader(att.Filename, bytes.NewReader(att.Content),
			mail.WithFileContentType(mail.ContentType(contentType))); err != nil {
			return err
		}
	}

	client, err := s.newClient()
	if err != nil {
		return err
	}
	return client.DialAndSendWithContext(ctx, msg)
}

// newClient construye el cliente SMTP según el puerto configurado.
func (s *EmailService) newClient() (*mail.Client, error) {
	opts := []mail.Option{
		mail.WithPort(s.port),
		mail.WithTimeout(emailDialTimeout),
	}

	if s.username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(s.username),
			mail.WithPassword(s.password),
		)
	}

	// Puerto 465 usa SSL/TLS implícito; el resto (ej: 587) usa STARTTLS.
	if s.port == 465 {
		opts = append(opts, mail.WithSSL())
	} else {
		opts = append(opts, mail.WithTLSPortPolicy(mail.TLSMandatory))
	}

	return mail.NewClient(s.host, opts...)
}
