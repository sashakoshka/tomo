package tomo

import "io"

// Data represents drag-and-drop, selection, or clipboard data.
type Data interface {
	io.Reader

	// Mime returns the MIME type of the data, such as text/plain,
	// text/html, image/png, etc.
	Mime () (mimeType Mime)

	// Convert attempts to convert the data to another MIME type. If the
	// data could not be converted, it should return an error.
	Convert (to Mime) (converted Data, err error)
}

// Mime represents a MIME type.
type Mime struct {
	// Type is the first half of the MIME type, and Subtype is the second
	// half. The separating slash is not included in either. For example,
	// text/html becomes:
	// Mime { Type: "text", Subtype: "html" }
	Type, Subtype string
}
