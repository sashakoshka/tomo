package ansi

// C0 represents a list of C0 control codes.
// https://en.wikipedia.org/wiki/C0_and_C1_control_codes
type C0 byte; const (
	C0_Null C0 = iota
	C0_StartOfHeading
	C0_StartOfText
	C0_EndOfText
	C0_EndOfTransmission
	C0_Enquiry
	C0_Acknowledge
	C0_Bell
	C0_Backspace
	C0_CharacterTab
	C0_LineFeed
	C0_LineTab
	C0_FormFeed
	C0_CarriageReturn
	C0_ShiftOut
	C0_ShiftIn
	C0_DataLinkEscape
	C0_DeviceControlOne
	C0_DeviceControlTwo
	C0_DeviceControlThree
	C0_DeviceControlFour
	C0_NegativeAcknowledge
	C0_SynchronousIdle
	C0_EndOfTransmissionBlock
	C0_Cancel
	C0_EndOfMedium
	C0_Substitute
	C0_Escape
	C0_FileSeparator
	C0_GroupSeparator
	C0_RecordSeparator
	C0_UnitSeparator
)
