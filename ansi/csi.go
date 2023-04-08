package ansi

// CSI represents a list of CSI sequences that have no parameters.
// FIXME: some of these do indeed have parameters
type CSI int; const (
	CSI_DeviceStatusReport CSI = iota
	CSI_SaveCursorPosition
	CSI_RestoreCursorPosition
	CSI_ShowCursor
	CSI_HideCursor
	CSI_EnableReportingFocus
	CSI_DisableReportingFocus
	CSI_EnableAlternativeBuffer
	CSI_DisableAlternativeBuffer
	CSI_EnableBracketedPasteMode
	CSI_DisableBracketedPasteMode
)
