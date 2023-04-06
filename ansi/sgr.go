package ansi

// SGR represents a list of Select Graphic Rendition parameters.
type SGR int; const (
	SGR_Normal SGR = iota
	SGR_Bold
	SGR_Dim
	SGR_Italic
	SGR_Underline
	SGR_SlowBlink
	SGR_RapidBlink
	SGR_Reverse
	SGR_Conceal
	SGR_Strike
	
	SGR_FontPrimary
	SGR_Font1
	SGR_Font2
	SGR_Font3
	SGR_Font4
	SGR_Font5
	SGR_Font6
	SGR_Font7
	SGR_Font8
	SGR_Font9
	SGR_FontFraktur
	
	SGR_DoubleUnderline
	SGR_NormalIntensity
	SGR_NeitherItalicNorBlackletter
	SGR_NotUnderlined
	SGR_NotBlinking
	SGR_PorportionalSpacing
	SGR_NotReversed
	SGR_NotCrossedOut
	
	SGR_ForegroundColorBlack
	SGR_ForegroundColorRed
	SGR_ForegroundColorGreen
	SGR_ForegroundColorYellow
	SGR_ForegroundColorBlue
	SGR_ForegroundColorMagenta
	SGR_ForegroundColorCyan
	SGR_ForegroundColorWhite
	SGR_ForegroundColor
	SGR_ForegroundColorDefault
	
	SGR_BackgroundColorBlack
	SGR_BackgroundColorRed
	SGR_BackgroundColorGreen
	SGR_BackgroundColorYellow
	SGR_BackgroundColorBlue
	SGR_BackgroundColorMagenta
	SGR_BackgroundColorCyan
	SGR_BackgroundColorWhite
	SGR_BackgroundColor
	SGR_BackgroundColorDefault

	SGR_DisablePorportionalSpacing
	SGR_Framed
	SGR_Encircled
	SGR_Overlined
	SGR_NeitherFramedNorEncircled
	SGR_NotOverlined
	
	SGR_UnderlineColor
	SGR_UnderlineColorDefault

	SGR_IdeogramUnderline
	SGR_IdeogramDoubleUnderline
	SGR_IdeogramOverline
	SGR_IdeogramDoubleOverline
	SGR_IdeogramStressMarking
	SGR_NoIdeogramAttributes

	SGR_Superscript
	SGR_Subscript
	SGR_NeitherSuperscriptNorSubscript

	SGR_ForegroundColorBrightBlack
	SGR_ForegroundColorBrightRed
	SGR_ForegroundColorBrightGreen
	SGR_ForegroundColorBrightYellow
	SGR_ForegroundColorBrightBlue
	SGR_ForegroundColorBrightMagenta
	SGR_ForegroundColorBrightCyan
	SGR_ForegroundColorBrightWhite
	
	SGR_BackgroundColorBrightBlack
	SGR_BackgroundColorBrightRed
	SGR_BackgroundColorBrightGreen
	SGR_BackgroundColorBrightYellow
	SGR_BackgroundColorBrightBlue
	SGR_BackgroundColorBrightMagenta
	SGR_BackgroundColorBrightCyan
	SGR_BackgroundColorBrightWhite
)
