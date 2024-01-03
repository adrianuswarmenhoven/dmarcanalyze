package main

import (
	"bytes"
	"encoding/base64"
	"strings"
)

func newAttachmentReaderFromMultipartMixed(attachment []byte) *bytes.Reader {
	// This is a bit messy, but it works
	multipartData := strings.SplitN(strings.TrimSpace(string(attachment)), "\n", 3) // Current delim, Content delim, body
	boundary := "--" + strings.TrimSpace(multipartData[1][strings.Index(multipartData[1], "boundary=")+9:])
	data := strings.Split(multipartData[2], boundary)
	AttachmentBody := ""
	attachmentType := ""
	for _, d := range data {
		if strings.Contains(d, "Content-Type: application/") {
			lines := strings.Split(d, "\n")
			headersDone := false
			attachmentType = ""
			for _, l := range lines {
				if !headersDone && strings.HasPrefix(l, "Content-Type: application/") {
					attachmentType = strings.TrimSpace(l[14:])
					if strings.Contains(attachmentType, ";") {
						attachmentType = strings.Split(attachmentType, ";")[0]
					}
					continue
				}
				if !headersDone && attachmentType != "" && strings.TrimSpace(l) == "" {
					headersDone = true
					continue
				}
				if headersDone {
					AttachmentBody += strings.TrimSpace(l)
				}
			}
		}
	}
	attachment = []byte(AttachmentBody)
	// Let's try and decode it. If it fails it is no problem
	base64.StdEncoding.Decode(attachment, attachment)
	return bytes.NewReader(attachment)
}
