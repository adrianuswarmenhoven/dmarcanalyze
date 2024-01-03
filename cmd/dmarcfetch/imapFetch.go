package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	report "github.com/oliverpool/go-dmarc-report"
)

// DialTLS connects to an IMAP server with implicit TLS.
func DialTLS(address string, insecure bool, options *imapclient.Options) (*imapclient.Client, error) {
	tlsconf := tls.Config{
		NextProtos: []string{"imap"},
	}
	if insecure {
		tlsconf.InsecureSkipVerify = true
	}
	conn, err := tls.Dial("tcp", address, &tlsconf)
	if err != nil {
		return nil, err
	}
	return imapclient.New(conn, options), nil
}

func getReportsViaIMAP4(server, user, password string, since, before time.Time) ([]*report.Aggregate, error) {
	slog.Debug("Connecting to IMAP4 server:", "server", server)
	client, err := DialTLS(server, true, nil)
	if err != nil {
		slog.Error("connect failed", "error", err)
		return nil, fmt.Errorf("connect failed: %w", err)
	}
	defer client.Close()

	slog.Debug("Logging in to IMAP4 server:", "user", user)
	cmd := client.Login(user, password)
	if err := cmd.Wait(); err != nil {
		slog.Error("login failed", "error", err)
		return nil, fmt.Errorf("login failed: %w", err)
	}

	slog.Debug("Selecting Agents.Dmarc")
	// select the mailbox
	if _, err := client.Select("Agents.Dmarc", &imap.SelectOptions{
		ReadOnly: true,
	}).Wait(); err != nil {
		slog.Error("select failed", "error", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	searchData, err := client.Search(&imap.SearchCriteria{
		Since:  since,
		Before: before,
		Header: []imap.SearchCriteriaHeaderField{
			{
				Key:   "Subject",
				Value: "Report Domain: ",
			},
		},
	}, nil).Wait()
	if err != nil {
		slog.Error("IMAP4 search failed", "error", err)
		return nil, fmt.Errorf("IMAP4 search failed: %w", err)
	}
	if len(searchData.All) == 0 {
		slog.Info("No reports found")
		return nil, fmt.Errorf("no reports found")
	}
	fetchOptions := &imap.FetchOptions{
		Flags:    true,
		Envelope: true,
		BodySection: []*imap.FetchItemBodySection{
			{
				Specifier: imap.PartSpecifierText,
			},
			{
				Specifier: imap.PartSpecifierHeader,
			},
		},
	}
	msgs, err := client.Fetch(searchData.All, fetchOptions).Collect()
	if err != nil {
		slog.Error("fetch failed", "error", err)
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	reports := make([]*report.Aggregate, 0)
	total := len(msgs)
	for id, msg := range msgs {
		slog.Debug("Fetching messages:", "current", id+1, "total", total) //id+1 because id starts at 0
		bodyTxt := []byte{}
		headersTxt := []byte{}
	bodyParts:
		for section, body := range msg.BodySection {
			switch section.Specifier {
			case imap.PartSpecifierText:
				bodyTxt = append(bodyTxt, body...)
			case imap.PartSpecifierHeader:
				headersTxt = body
			default:
				slog.Debug("Unknown body section", "section", section, "body", body)
				continue bodyParts
			}
		}
		messageToParse := string(headersTxt) + string(bodyTxt)
		email, err := Parse(bytes.NewReader([]byte(messageToParse)))
		if err != nil {
			// handle error
			slog.Error("parse failed", "error", err)
		}
		attachment := []byte{}
		attachmentType := ""
		if len(email.Attachments) < 1 { // Empty message, only data (Google does this)
			attachmentType = email.Header.Get("Content-Type")
			attachment = bodyTxt
		} else {
			for _, a := range email.Attachments {
				attachmentBytes, err := io.ReadAll(a.Data)
				if err != nil {
					slog.Error("attachment read failed", "error", err)
				}
				if strings.HasPrefix(a.ContentType, "text/plain") { // FIXME, maybe uncompressed?
					continue
				}
				attachment = attachmentBytes
				attachmentType = a.ContentType
			}
		}
		if len(attachment) < 1 || attachmentType == "" {
			slog.Error("No attachment found", "message", messageToParse)
			return nil, fmt.Errorf("no attachment found")
		}

		// Let's try and decode it. If it fails it is no problem
		base64.StdEncoding.Decode(attachment, attachment)

		attachmentReader := bytes.NewReader(attachment)
		agg := &report.Aggregate{}
		if strings.Index(attachmentType, ";") > 0 {
			attachmentType = strings.Split(attachmentType, ";")[0]
		}
		err = nil
		switch attachmentType {
		case "application/zip":
			agg, err = report.DecodeZip(attachmentReader, int64(len(attachment)))
		case "application/x-gzip-compressed", "application/gzip":
			agg, err = report.DecodeGzip(attachmentReader)
		case "text/plain":
			agg, err = report.Decode(attachmentReader)
		case "application/x-zip-compressed":
			agg, err = report.DecodeZip(attachmentReader, int64(len(attachment)))
		case "multipart/mixed":
			attachmentReader = newAttachmentReaderFromMultipartMixed(attachment)
			agg, err = report.DecodeGzip(attachmentReader)
		default:
			slog.Debug("unknown type", "type", attachmentType)
			err = fmt.Errorf("unknown type '%s'", attachmentType)
		}
		if err != nil {
			slog.Error("decode failed", "error", err, "attachment", string(attachment))
			return nil, fmt.Errorf("decode failed: %w", err)
		}
		//fmt.Println(fmt.Sprintf("%#v", agg))
		reports = append(reports, agg)
	}
	return reports, nil
}
