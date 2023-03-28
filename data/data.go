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

// M is shorthand for creating a MIME type.
func M (ty, subtype string) Mime {
	return Mime { ty, subtype }
}

// String returns the string representation of the MIME type.
func (mime Mime) String () string {
	return mime.Type + "/" + mime.Subtype
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
