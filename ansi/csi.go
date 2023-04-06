package ansi

// CSI represents a list of CSI sequences that have no parameters.
type CSI int; const (
	CSI_AuxPortOn CSI = iota
	CSI_AuxPortOff
	CSI_DeviceStatusReport
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
