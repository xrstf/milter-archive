package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/phalaaxx/milter"
)

type archiveMilter struct {
	target   string
	spamMode bool
	decoder  *mime.WordDecoder

	message *bytes.Buffer
	subject string
	isSpam  bool
}

func newArchiveMilter(target string, spamMode bool) *archiveMilter {
	decoder := new(mime.WordDecoder)

	return &archiveMilter{
		target:   target,
		spamMode: spamMode,
		decoder:  decoder,
	}
}

// Connect is called to provide SMTP connection data for incoming message
func (m *archiveMilter) Connect(host string, family string, port uint16, addr net.IP, mod *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// Helo is called to process any HELO/EHLO related filters
func (m *archiveMilter) Helo(name string, mod *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// MailFrom is called to process filters on envelope FROM address
func (m *archiveMilter) MailFrom(from string, mod *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// RcptTo is called to process filters on envelope TO address
func (m *archiveMilter) RcptTo(rcptTo string, mod *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// Header is called once for each header in incoming message
func (m *archiveMilter) Header(name string, value string, mod *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

// Headers is called when all message headers have been processed
func (m *archiveMilter) Headers(h textproto.MIMEHeader, mod *milter.Modifier) (milter.Response, error) {
	// prepare message buffer
	m.message = new(bytes.Buffer)

	// print headers to message buffer
	for k, vl := range h {
		for _, v := range vl {
			if _, err := fmt.Fprintf(m.message, "%s: %s\n", k, v); err != nil {
				return nil, err
			}
		}
	}

	if _, err := fmt.Fprintf(m.message, "\n"); err != nil {
		return nil, err
	}

	m.subject = h.Get("Subject")
	m.isSpam = strings.ToLower(h.Get("X-Spam")) == "yes"

	return milter.RespContinue, nil
}

// BodyChunk is called to process next message body chunk data (up to 64KB in size)
func (m *archiveMilter) BodyChunk(chunk []byte, mod *milter.Modifier) (milter.Response, error) {
	if _, err := m.message.Write(chunk); err != nil {
		return nil, err
	}

	return milter.RespContinue, nil
}

// Body is called at the end of each message
func (m *archiveMilter) Body(mod *milter.Modifier) (milter.Response, error) {
	ts := time.Now().UTC().Format("2006-01-02T150405.000Z")

	decoded, err := m.decoder.DecodeHeader(m.subject)
	if err != nil {
		log.Printf("Error decoding subject: %v", err)
	} else {
		m.subject = decoded
	}

	slug := slug.Make(m.subject)
	if len(slug) == 0 {
		slug = "no-subject"
	}

	directory := m.target

	if m.spamMode {
		if m.isSpam {
			directory = filepath.Join(directory, "spam")
		} else {
			directory = filepath.Join(directory, "ham")
		}

		if _, err := os.Stat(directory); err != nil {
			if err := os.Mkdir(directory, 0755); err != nil {
				log.Printf("Failed to create spam subdirectory %s: %v", directory, err)
				return milter.RespContinue, nil
			}
		}
	}

	filename := filepath.Join(directory, fmt.Sprintf("%s-%s.eml", ts, slug))

	err = ioutil.WriteFile(filename, m.message.Bytes(), 0644)
	if err != nil {
		log.Printf("Error writing e-mail to %s: %v", filename, err)
	}

	return milter.RespContinue, nil
}
