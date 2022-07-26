package tmpmail

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/mhale/smtpd"
	"go.uber.org/zap"

	"tmpmail/email"
	"tmpmail/entity"
)

type SMTPServerStorage interface {
	AccountExists(username string) (bool, error)
	AddEmail(username string, email entity.Email) error
}

type SMTPServer struct {
	server *smtpd.Server
}

func NewSMTPServer(l *zap.Logger, s SMTPServerStorage, tc *tls.Config, addr, domain, mailDomain string) *SMTPServer {

	smtpd.Debug = true

	srv := &smtpd.Server{
		Addr:      addr,
		Appname:   "tmpmail",
		TLSConfig: tc,
		Handler: func(remoteAddr net.Addr, from string, to []string, data []byte) error {
			m, err := email.Parse(bytes.NewReader(data))
			if err != nil {
				return fmt.Errorf("parse email: %w", err)
			}

			mm := entity.Email{
				Subject:         m.Subject,
				Date:            m.Date,
				MessageID:       m.MessageID,
				InReplyTo:       m.InReplyTo,
				References:      m.References,
				ResentMessageID: m.ResentMessageID,
				ContentType:     m.ContentType,
				HTMLBody:        m.HTMLBody,
				TextBody:        m.TextBody,
			}

			if mm.Date.IsZero() {
				mm.Date = time.Now()
			}
			if m.Sender != nil {
				mm.Sender = m.Sender.String()
			}
			if m.ResentSender != nil {
				mm.ResentSender = m.ResentSender.String()
			}
			for _, i := range m.From {
				mm.From = append(mm.From, i.String())
			}
			for _, i := range m.ResentTo {
				mm.ResentTo = append(mm.ResentTo, i.String())
			}
			if !m.ResentDate.IsZero() {
				mm.ResentDate = &m.ResentDate
			}
			for _, i := range m.To {
				mm.To = append(mm.To, i.String())
			}
			for _, i := range m.Cc {
				mm.Cc = append(mm.Cc, i.String())
			}
			for _, i := range m.Bcc {
				mm.Bcc = append(mm.Bcc, i.String())
			}
			for _, i := range m.ResentFrom {
				mm.ResentFrom = append(mm.ResentFrom, i.String())
			}
			for _, i := range m.ResentTo {
				mm.ResentTo = append(mm.ResentTo, i.String())
			}
			for _, i := range m.ResentCc {
				mm.ResentCc = append(mm.ResentCc, i.String())
			}
			for _, i := range m.ResentBcc {
				mm.ResentBcc = append(mm.ResentBcc, i.String())
			}
			for _, a := range m.Attachments {
				data, err := io.ReadAll(a.Data)
				if err != nil {
					return fmt.Errorf("read attachment: %w", err)
				}
				mm.Attachments = append(mm.Attachments, entity.Attachment{
					Filename:    a.Filename,
					ContentType: a.ContentType,
					Data:        data,
				})
			}
			for _, a := range m.EmbeddedFiles {
				data, err := io.ReadAll(a.Data)
				if err != nil {
					return fmt.Errorf("read embedded file: %w", err)
				}
				mm.EmbeddedFiles = append(mm.EmbeddedFiles, entity.EmbeddedFile{
					CID:         a.CID,
					ContentType: a.ContentType,
					Data:        data,
				})
			}

			logger := l.With(
				zap.String("remote_addr", remoteAddr.String()),
				zap.Strings("from", mm.From))

			if len(mm.To) == 0 {
				logger.Warn("email without recipients")
				return nil
			}

			var wg sync.WaitGroup
			for i := range m.To {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

					username := strings.TrimSuffix(m.To[i].Address, "@"+domain)
					if len(username) == len(m.To[i].Address) {
						return
					}

					logger = logger.With(
						zap.String("to", mm.To[i]),
						zap.String("username", username))

					err = s.AddEmail(username, mm)
					if err != nil {
						logger.Error("add email data to storage", zap.Error(err))
						return
					}

					logger.Info("email data added")
				}(i)
			}
			wg.Wait()
			return nil
		},
		HandlerRcpt: func(remoteAddr net.Addr, from string, to string) bool {
			username := strings.TrimSuffix(to, "@"+domain)
			if len(username) == len(to) {
				return false
			}
			exist, err := s.AccountExists(username)
			if err != nil {
				l.Error("check account exists", zap.String("account", username))
				return false
			}
			return exist
		},
		Hostname: mailDomain,
		LogRead: func(remoteIP, verb, line string) {
			l.Debug("smtp read",
				zap.String("remote_ip", remoteIP),
				zap.String("line", line))
		},
		LogWrite: func(remoteIP, verb, line string) {
			l.Debug("smtp write",
				zap.String("remote_ip", remoteIP),
				zap.String("line", line))
		},
	}

	return &SMTPServer{server: srv}
}

func (s *SMTPServer) ListenAndServe() error {
	err := s.server.ListenAndServe()
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return fmt.Errorf("listen and server: %w", err)
	}
	return nil
}

func (s *SMTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
