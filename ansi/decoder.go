package ansi

import "image/color"

// Useful: https://en.wikipedia.org/wiki/ANSI_escape_code

// Decoder is a state machine capable of decoding text contianing various escape
// codes and sequences. It satisfies io.Writer. It has no constructor and its
// zero value can be used safely.
type Decoder struct {
	// OnText is called when a segment of text is processed.
	OnText func (string)

	// OnC0 is called when a C0 control code is processed that isn't a
	// whitespace character.
	OnC0 func (C0)

	// OnC1 is called when a C1 escape sequence is processed.
	OnC1 func (C1)

	// OnCSI is called when a non-SGR CSI escape sequence with no parameters
	// is processed.
	OnCSI func (CSI)

	// OnSGR is called when a CSI SGR escape sequence with no parameters is
	// processed.
	OnSGR func (SGR)

	// Non-SGR CSI sequences with parameters:
	OnCursorUp                   func (distance int)
	OnCursorDown                 func (distance int)
	OnCursorForward              func (distance int)
	OnCursorBack                 func (distance int)
	OnCursorNextLine             func (distance int)
	OnCursorPreviousLine         func (distance int)
	OnCursorHorizontalAbsolute   func (column int)
	OnCursorPosition             func (column, row int)
	OnEraseInDisplay             func (mode int)
	OnEraseInLine                func (mode int)
	OnScrollUp                   func (distance int)
	OnScrollDown                 func (distance int)
	OnHorizontalVerticalPosition func (column, row int)

	// SGR CSI sequences with parameters:
	OnForegroundColor     func (Color)
	OnForegroundColorTrue func (color.RGBA)
	OnBackgroundColor     func (Color)
	OnBackgroundColorTrue func (color.RGBA)
	OnUnderlineColor      func (Color)
	OnUnderlineColorTrue  func (color.RGBA)

	// OSC sequences from XTerm:
	OnWindowTitle     func (title string)
	OnIconName        func (name string)
	OnIconFile        func (path string)
	OnXProperty       func (property, value string)
	OnSelectionPut    func (selection, text string)
	OnSelectionGet    func (selection string)
	OnQueryAllowed    func ()
	OnQueryDisallowed func ()

	// OSC sequences from iTerm2:
	OnCursorShape     func (shape int)
	OnHyperlink       func (params map[string] string, text, link string)
	OnBackgroundImage func (path string)

	state decodeState
	csiParameter []byte
	csiIdentifier byte
}

type decodeState int; const (
	decodeStateText decodeState = iota
	decodeStateAwaitC1
	decodeStateGatherCSI
)

func (decoder *Decoder) Write (buffer []byte) (wrote int, err error) {
	wrote = len(buffer)

	for len(buffer) > 0 {
		switch decoder.state {
		case decodeStateText:
			if buffer[0] == byte(C0_Escape) {
				// begin C1 control code
				decoder.state = decodeStateAwaitC1
				buffer = buffer[1:]
				
			} else if buffer[0] < ' ' {
				// process C0 control code
				if decoder.OnC0 != nil {
					decoder.OnC0(C0(buffer[0]))
				}
				buffer = buffer[1:]
				
			} else {
				// process as much plain text as we can
				buffer = decoder.processString(buffer)
			}
			
		case decodeStateAwaitC1:
			// TODO: handle OSC sequences
			// TODO: handle device control string
			// TODO: handle privacy message
			// TODO: handle application program command
		
			if buffer[0] < 128 {
				// false alarm, this is just a C0 escape
				if decoder.OnC0 != nil {
					decoder.OnC0(C0_Escape)
				}
				
			} else if buffer[0] == byte(C1_ControlSequenceIntroducer) {
				// abandon all hope ye who enter here
				decoder.state = decodeStateGatherCSI
				decoder.csiParameter  = nil
				decoder.csiIdentifier = 0
				
			} else {
				// process C1 control code
				if decoder.OnC1 != nil {
					decoder.OnC1(C1(buffer[0]))
				}
			}
			buffer = buffer[1:]

		case decodeStateGatherCSI:
			if buffer[0] < 0x30 || buffer[0] > 0x3F {
				decoder.csiIdentifier = buffer[0]
				decoder.processCSI()
			} else {
				decoder.csiParameter = append (
					decoder.csiParameter,
					buffer[0])
			}
			buffer = buffer[1:]
		}
	}
	
	return
}

func (decoder *Decoder) processString (buffer []byte) []byte {
	for index, char := range buffer {
		if C0_Escape.Is(char) {
			if decoder.OnText != nil {
				decoder.OnText(string(buffer[:index]))
			}
			return buffer[:index]
		}
	}
	return buffer
}

func (decoder *Decoder) processCSI () {
	// TODO: analyze CSI parameter and id
}
