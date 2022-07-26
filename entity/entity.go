package entity

import "time"

type Attachment struct {
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Data        []byte `json:"data,omitempty"`
}

type EmbeddedFile struct {
	CID         string `json:"cid,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Data        []byte `json:"data,omitempty"`
}

type Email struct {
	Subject    string    `json:"subject,omitempty"`
	Sender     string    `json:"sender,omitempty"`
	From       []string  `json:"from,omitempty"`
	ReplyTo    []string  `json:"replyTo,omitempty"`
	To         []string  `json:"to,omitempty"`
	Cc         []string  `json:"cc,omitempty"`
	Bcc        []string  `json:"bcc,omitempty"`
	Date       time.Time `json:"date,omitempty"`
	MessageID  string    `json:"messageID,omitempty"`
	InReplyTo  []string  `json:"inReplyTo,omitempty"`
	References []string  `json:"references,omitempty"`

	ResentFrom      []string   `json:"resentFrom,omitempty"`
	ResentSender    string     `json:"resentSender,omitempty"`
	ResentTo        []string   `json:"resentTo,omitempty"`
	ResentDate      *time.Time `json:"resentDate,omitempty"`
	ResentCc        []string   `json:"resentCc,omitempty"`
	ResentBcc       []string   `json:"resentBcc,omitempty"`
	ResentMessageID string     `json:"resentMessageID,omitempty"`

	ContentType string `json:"contentType,omitempty"`

	HTMLBody string `json:"htmlBody,omitempty"`
	TextBody string `json:"textBody,omitempty"`

	Attachments   []Attachment   `json:"attachments,omitempty"`
	EmbeddedFiles []EmbeddedFile `json:"embeddedFiles,omitempty"`
}

type Account struct {
	Username string  `json:"username"`
	TTL      int64   `json:"ttl"`
	Emails   []Email `json:"emails"`
}
