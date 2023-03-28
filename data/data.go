package data

import "io"
import "bytes"

// Data represents arbitrary polymorphic data that can be used for data transfer
// between applications.
type Data map[Mime] io.ReadCloser

// Mime represents a MIME type.
type Mime struct {
	// Type is the first half of the MIME type, and Subtype is the second
	// half. The separating slash is not included in either. For example,
	// text/html becomes:
	// Mime { Type: "text", Subtype: "html" }
	Type, Subtype string
}

var MimePlain = Mime { "text", "plain" }

var MimeFile = Mime { "text", "uri-list" }

type byteReadCloser struct { *bytes.Reader }
func (byteReadCloser) Close () error { return nil }

// Text returns plain text Data given a string.
func Text (text string) Data {
	return Data {
		MimePlain: byteReadCloser { bytes.NewReader([]byte(text)) },
	}
}
