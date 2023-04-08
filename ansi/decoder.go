package ansi

import "strconv"
import "strings"
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

	// OnDCS is called when a device control string is processed.
	OnDCS func (string)

	// OnCSI is called when a non-SGR CSI escape sequence with no parameters
	// is processed.
	OnCSI func (CSI)

	// OnSGR is called when a CSI SGR escape sequence with no parameters is
	// processed.
	OnSGR func (SGR)

	// OnPM is called when a privacy message is processed.
	OnPM func (string)

	// OnAPC is called when an application program command is processed.
	OnAPC func (string)

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
	expectingST bool
	gathered []byte
}

type decodeState int; const (
	decodeStateText decodeState = iota
	decodeStateAwaitC1
	
	decodeStateGatherDCS
	decodeStateGatherSOS
	decodeStateGatherCSI
	decodeStateGatherOSC
	decodeStateGatherPM
	decodeStateGatherAPC
)

func (decoder *Decoder) Write (buffer []byte) (wrote int, err error) {
	wrote = len(buffer)

	for len(buffer) > 0 {
		switch decoder.state {
		case decodeStateText:
			if C0_Escape.Is(buffer[0]) {
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
			if buffer[0] < 128 {
				// false alarm, this is just a C0 escape
				if decoder.OnC0 != nil {
					decoder.OnC0(C0_Escape)
				}
				break
			}

			switch C1(buffer[0]) {
			case C1_DeviceControlString:
				decoder.state = decodeStateGatherDCS
			case C1_StartOfString:
				decoder.state = decodeStateGatherSOS
			case C1_ControlSequenceIntroducer:
				decoder.state = decodeStateGatherCSI
			case C1_OperatingSystemCommand:
				decoder.state = decodeStateGatherOSC
			case C1_PrivacyMessage:
				decoder.state = decodeStateGatherPM
			case C1_ApplicationProgramCommand:
				decoder.state = decodeStateGatherAPC
			default:
				// process C1 control code
				if decoder.OnC1 != nil {
					decoder.OnC1(C1(buffer[0]))
				}
			}
			buffer = buffer[1:]

		case 
			decodeStateGatherDCS,
			decodeStateGatherSOS,
			decodeStateGatherOSC,
			decodeStateGatherPM,
			decodeStateGatherAPC:

			if decoder.expectingST && C1_StringTerminator.Is(buffer[0]) {
				// remove the trailing ESC
				decoder.gathered = decoder.gathered [
					:len(decoder.gathered) - 1]
				
				if decoder.state == decodeStateGatherOSC {
					// we understand some OSC codes so we
					// handle them differently
					decoder.processOSC()
				} else {
					// generic handler for uncommon stuff
					decoder.processGeneric()
				}
				decoder.state = decodeStateText
			}
			if C0_Escape.Is(buffer[0]) {
				decoder.expectingST = true
			}
			
		case decodeStateGatherCSI:
			decoder.gather(buffer[0])
			if buffer[0] < 0x30 || buffer[0] > 0x3F {
				decoder.processCSI()
				decoder.state = decodeStateText
			}
			buffer = buffer[1:]
		}
	}
	
	return
}

func (decoder *Decoder) gather (character byte) {
	decoder.gathered = append(decoder.gathered, character)
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

func (decoder *Decoder) processGeneric () {
	parameter := string(decoder.gathered)
	switch decoder.state {
	case decodeStateGatherDCS:
		if decoder.OnDCS  != nil { decoder.OnDCS(parameter)  }
	case decodeStateGatherSOS:
		if decoder.OnText != nil { decoder.OnText(parameter) }
	case decodeStateGatherPM:
		if decoder.OnPM   != nil { decoder.OnPM(parameter)   }
	case decodeStateGatherAPC:
		if decoder.OnAPC  != nil { decoder.OnAPC(parameter)  }
	}
}

func (decoder *Decoder) processOSC () {
	// TODO: analyze OSC
}

func (decoder *Decoder) processCSI () {
	if len(decoder.gathered) < 2 { return }
	parameters := ParameterInts(decoder.gathered)

	var p0, p1, p2, p3 int
	if len(parameters) > 0 { p0 = parameters[0] }
	if len(parameters) > 1 { p1 = parameters[1] }
	if len(parameters) > 2 { p2 = parameters[2] }
	if len(parameters) > 3 { p3 = parameters[3] }
	
	switch Last(decoder.gathered) {
	case 'A': if decoder.OnCursorUp                 != nil { decoder.OnCursorUp(clampOne(p0)) }
	case 'B': if decoder.OnCursorDown               != nil { decoder.OnCursorDown(clampOne(p0)) }
	case 'C': if decoder.OnCursorForward            != nil { decoder.OnCursorForward(clampOne(p0)) }
	case 'D': if decoder.OnCursorBack               != nil { decoder.OnCursorBack(clampOne(p0)) }
	case 'E': if decoder.OnCursorNextLine           != nil { decoder.OnCursorNextLine(clampOne(p0)) }
	case 'F': if decoder.OnCursorPreviousLine       != nil { decoder.OnCursorPreviousLine(clampOne(p0)) }
	case 'G': if decoder.OnCursorHorizontalAbsolute != nil { decoder.OnCursorHorizontalAbsolute(clampOne(p0)) }
	case 'H', 'f': if decoder.OnCursorPosition      != nil { decoder.OnCursorPosition(clampOne(p0), clampOne(p1)) }
	
	case 'J': if decoder.OnEraseInDisplay != nil { decoder.OnEraseInDisplay(p0) }
	case 'K': if decoder.OnEraseInLine    != nil { decoder.OnEraseInLine(p0) }
	case 'S': if decoder.OnScrollUp       != nil { decoder.OnScrollUp(clampOne(p0)) }
	case 'T': if decoder.OnScrollDown     != nil { decoder.OnScrollDown(clampOne(p0)) }
	
	case 'm':
		p0 := SGR(p0)
		switch {
		case
			p0 >= SGR_ForegroundColorBlack &&
			p0 <= SGR_ForegroundColorWhite &&
			decoder.OnForegroundColor != nil :
			decoder.OnForegroundColor(Color(p0 - SGR_ForegroundColorBlack))
		case
			p0 >= SGR_ForegroundColorBrightBlack &&
			p0 <= SGR_ForegroundColorBrightWhite &&
			decoder.OnForegroundColor != nil :
			decoder.OnForegroundColor(Color(p0 - SGR_ForegroundColorBrightBlack + 8))
			
		case p0 == SGR_ForegroundColor:
			switch p1 {
			case 2:
				if decoder.OnForegroundColor == nil { break }
				decoder.OnForegroundColor(Color(p1))
			case 5:
				if decoder.OnForegroundColorTrue == nil { break }
				decoder.OnForegroundColorTrue (color.RGBA {
					R: uint8(p1),
					G: uint8(p2),
					B: uint8(p3),
					A: 0xFF,
				})
			}
		
		case
			p0 >= SGR_BackgroundColorBlack &&
			p0 <= SGR_BackgroundColorWhite &&
			decoder.OnBackgroundColor != nil :
			decoder.OnBackgroundColor(Color(p0 - SGR_BackgroundColorBlack))
		case
			p0 >= SGR_BackgroundColorBrightBlack &&
			p0 <= SGR_BackgroundColorBrightWhite &&
			decoder.OnBackgroundColor != nil :
			decoder.OnBackgroundColor(Color(p0 - SGR_BackgroundColorBrightBlack + 8))
			
		case p0 == SGR_BackgroundColor:
			switch p1 {
			case 2:
				if decoder.OnBackgroundColor == nil { break }
				decoder.OnBackgroundColor(Color(p1))
			case 5:
				if decoder.OnBackgroundColorTrue == nil { break }
				decoder.OnBackgroundColorTrue (color.RGBA {
					R: uint8(p1),
					G: uint8(p2),
					B: uint8(p3),
					A: 0xFF,
				})
			}
		
		case p0 == SGR_UnderlineColor:
			switch p1 {
			case 2:
				if decoder.OnUnderlineColor == nil { break }
				decoder.OnUnderlineColor(Color(p1))
			case 5:
				if decoder.OnUnderlineColorTrue == nil { break }
				decoder.OnUnderlineColorTrue (color.RGBA {
					R: uint8(p1),
					G: uint8(p2),
					B: uint8(p3),
					A: 0xFF,
				})
			}

		default: if decoder.OnSGR != nil { decoder.OnSGR(SGR(p0)) }
		}

	// TODO
	case 'n':
	case 's':
	case 'u':
	case 'h':
	case 'l':
	}
}

func clampOne (number int) int {
	if number < 1 {
		return 1
	} else {
		return number
	}
}

// Last returns the last item of a slice.
func Last[T any] (source []T) T {
	return source[len(source) - 1]
}

// ParameterStrings separates a byte slice by semicolons into a list of strings.
func ParameterStrings (source []byte) (parameters []string) {
	parameters = strings.Split(string(source), ";")
	for index := range parameters {
		parameters[index] = strings.TrimSpace(parameters[index])
	}
	return
}

// ParameterInts is like ParameterStrings, but returns integers instead of
// strings. If a parameter is empty or cannot be converted into an integer, that
// parameter will be zero.
func ParameterInts (source []byte) (parameters []int) {
	stringParameters := ParameterStrings(source)
	parameters = make([]int, len(stringParameters))
	for index, parameter := range stringParameters {
		parameters[index], _ = strconv.Atoi(parameter)
	}
	return
}
