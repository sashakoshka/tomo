package ansi

// C1 represents a list of C1 control codes.
// https://en.wikipedia.org/wiki/C0_and_C1_control_codes
type C1 byte; const (
	C1_PaddingCharacter C1 = iota + 128
	C1_HighOctetPreset
	C1_BreakPermittedHere
	C1_NoBreakHere
	C1_Index
	C1_NextLine
	C1_StartOfSelectedArea
	C1_EndOfSelectedArea
	C1_CharacterTabSet
	C1_CharacterTabWithJustification
	C1_LineTabSet
	C1_PartialLineForward
	C1_PartialLineBackward
	C1_ReverseLineFeed
	C1_SingleShift2
	C1_SingleShift3
	C1_DeviceControlString
	C1_PrivateUse1
	C1_PrivateUse2
	C1_SetTransmitState
	C1_CancelCharacter
	C1_MessageWaiting
	C1_StartOfProtectedArea
	C1_EndOfProtectedArea
	C1_StartOfString
	C1_SingleGraphicCharacterIntroducer
	C1_SingleCharacterIntroducer
	C1_ControlSequenceIntroducer
	C1_StringTerminator
	C1_OperatingSystemCommand
	C1_PrivacyMessage
	C1_ApplicationProgramCommand
)

// Is checks if a byte is equal to a C0 code.
func (code C0) Is (test byte) bool {
	return test == byte(code)
}

// Is checks if a byte is equal to a C1 code.
func (code C1) Is (test byte) bool {
	return byte(code) == test
}
