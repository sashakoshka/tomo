package tomo

import "image"
import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Color lits a number of cannonical colors, each with its own ID.
type Color int; const (
	// The sixteen ANSI terminal colors:
	ColorBlack Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorPurple
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightPurple
	ColorBrightCyan
	ColorBrightWhite

	// ColorForeground is the text/icon color of the theme.
	ColorForeground

	// ColorMidground is a generic raised element color.
	ColorMidground

	// ColorBackground is the background color of the theme.
	ColorBackground

	// ColorShadow is a generic shadow color.
	ColorShadow

	// ColorShine is a generic highlight color.
	ColorShine

	// ColorAccent is the accent color of the theme.
	ColorAccent
)

// Pattern lists a number of cannonical pattern types, each with its own ID.
type Pattern int; const (
	// PatternBackground is the window background of the theme. It appears
	// in things like containers and behind text.
	PatternBackground Pattern = iota

	// PatternDead is a pattern that is displayed on a "dead area" where no
	// controls exist, but there still must be some indication of visual
	// structure (such as in the corner between two scroll bars).
	PatternDead

	// PatternRaised is a generic raised pattern.
	PatternRaised

	// PatternSunken is a generic sunken pattern.
	PatternSunken

	// PatternPinboard is similar to PatternSunken, but it is textured.
	PatternPinboard

	// PatternButton is a button pattern.
	PatternButton

	// PatternInput is a pattern for input fields, editable text areas, etc.
	PatternInput

	// PatternGutter is a track for things to slide on.
	PatternGutter

	// PatternHandle is a handle that slides along a gutter.
	PatternHandle

	// PatternLine is an engraved line that separates things.
	PatternLine

	// PatternMercury is a fill pattern for progress bars, meters, etc.
	PatternMercury

	// PatternTableHead is a table row or column heading background.
	PatternTableHead

	// PatternTableCell is a table cell background.
	PatternTableCell

	// PatternLamp is an indicator light pattern.
	PatternLamp
)

// IconSize is a type representing valid icon sizes.
type IconSize int

const (
	IconSizeSmall IconSize = 16
	IconSizeLarge IconSize = 48
)

// Icon lists a number of cannonical icons, each with its own ID.
type Icon int

// IconNone specifies no icon.
const IconNone = -1

const (
	// Place icons
	IconHome Icon = iota
	Icon3DObjects
	IconPictures
	IconVideos
	IconMusic
	IconArchives
	IconBooks
	IconDocuments
	IconFonts
	IconPrograms
	IconLibraries
	IconDownloads
	IconRepositories
	IconSettings
	IconHistory)

const (
	// Object icons
	IconFile Icon = iota + 0x80
	IconDirectory
	IconPopulatedDirectory
	
	IconStorage
	IconMagneticTape
	IconFloppyDisk
	IconHDD
	IconSSD
	IconFlashDrive
	IconMemoryCard
	IconRomDisk
	IconRamDisk
	IconCD
	IconDVD

	IconNetwork
	IconInternet

	IconDevice
	IconServer
	IconNetworkSwitch
	IconRouter
	IconDesktop
	IconLaptop
	IconTablet
	IconPhone
	IconCamera

	IconPeripheral
	IconKeyboard
	IconMouse
	IconTrackpad
	IconPenTablet
	IconMonitor
	IconSpeaker
	IconMicrophone
	IconWebcam
	IconGameController

	IconPort
	IconNetworkPort
	IconUSBPort
	IconParallelPort
	IconSerialPort
	IconPS2Port
	IconMonitorPort)

const (
	// Action icons
	IconOpen Icon = iota + 0x100
	IconSave
	IconSaveAs
	IconNew
	IconNewFolder
	IconDelete

	IconCut
	IconCopy
	IconPaste

	IconAdd
	IconRemove
	IconAddBookmark
	IconRemoveBookmark
	IconAddFavorite
	IconRemoveFavorite
	
	IconPlay
	IconPause
	IconStop
	IconFastForward
	IconRewind
	IconToEnd
	IconToBeginning
	IconRecord
	IconVolumeUp
	IconVolumeDown
	IconMute

	IconBackward
	IconForward
	IconUpward
	IconRefresh

	IconYes
	IconNo

	IconUndo
	IconRedo

	IconRun
	IconSearch

	IconClose
	IconQuit
	IconIconify
	IconShade
	IconMaximize
	IconRestore

	IconReplace
	IconUnite
	IconDiffer
	IconInvert
	IconIntersect

	IconExpand)

const (
	// Status icons
	IconInformation Icon = iota + 0x180
	IconQuestion
	IconWarning
	IconError)

const (
	// Tool icons
	IconCursor Icon = iota + 0x200
	IconMeasure
	
	IconSelect
	IconSelectRectangle
	IconSelectEllipse
	IconSelectGeometric
	IconSelectFreeform
	IconSelectLasso
	IconSelectFuzzy
	
	IconTransform
	IconTranslate
	IconRotate
	IconScale
	IconWarp
	IconDistort
	
	IconPencil
	IconBrush
	IconEraser
	IconFill
	IconText)
	
// FontSize specifies the general size of a font face in a semantic way.
type FontSize int; const (
	// FontSizeNormal is the default font size that should be used for most
	// things.
	FontSizeNormal FontSize = iota

	// FontSizeLarge is a larger font size suitable for things like section
	// headings.
	FontSizeLarge

	// FontSizeHuge is a very large font size suitable for things like
	// titles, wizard step names, digital clocks, etc.
	FontSizeHuge

	// FontSizeSmall is a smaller font size. Try not to use this unless it
	// makes a lot of sense to do so, because it can negatively impact
	// accessibility. It is useful for things like copyright notices at the
	// bottom of some window that the average user doesn't actually care
	// about.
	FontSizeSmall
)

// FontStyle specifies stylistic alterations to a font face.
type FontStyle int; const (
	FontStyleRegular    FontStyle = 0
	FontStyleBold       FontStyle = 1
	FontStyleItalic     FontStyle = 2
	FontStyleMonospace  FontStyle = 4
	FontStyleBoldItalic FontStyle = 1 | 2
)

// Hints specifies rendering hints for a particular pattern. Elements can take
// these into account in order to gain extra performance.
type Hints struct {
	// StaticInset defines an inset rectangular area in the middle of the
	// pattern that does not change between PatternStates. If the inset is
	// zero on all sides, this hint does not apply.
	StaticInset artist.Inset

	// Uniform specifies a singular color for the entire pattern. If the
	// alpha channel is zero, this hint does not apply.
	Uniform color.RGBA
}

// Theme represents a visual style configuration,
type Theme interface {
	// FontFace returns the proper font for a given style, size, and case.
	FontFace (FontStyle, FontSize, Case) font.Face

	// Icon returns an appropriate icon given an icon name, size, and case.
	Icon (Icon, IconSize, Case) artist.Icon
	
	// Icon returns an appropriate icon given a file mime type, size, and,
	// case.
	MimeIcon (data.Mime, IconSize, Case) artist.Icon

	// Pattern returns an appropriate pattern given a pattern name, case,
	// and state.
	Pattern (Pattern, State, Case) artist.Pattern

	// Color returns an appropriate pattern given a color name, case, and
	// state.
	Color (Color, State, Case) color.RGBA

	// Padding returns how much space should be between the bounds of a
	// pattern whatever an element draws inside of it.
	Padding (Pattern, Case) artist.Inset

	// Margin returns the left/right (x) and top/bottom (y) margins that
	// should be put between any self-contained objects drawn within this
	// pattern (if applicable).
	Margin (Pattern, Case) image.Point

	// Sink returns a vector that should be added to an element's inner
	// content when it is pressed down (if applicable) to simulate a 3D
	// sinking effect.
	Sink (Pattern, Case) image.Point

	// Hints returns rendering optimization hints for a particular pattern.
	// These are optional, but following them may result in improved
	// performance.
	Hints (Pattern, Case) Hints
}

// Case sepecifies what kind of element is using a pattern. It contains a
// namespace parameter, an element parameter, and an optional component trail.
// All parameter values should be written in camel case. Themes can change their
// styling based on the case for fine-grained control over the look and feel of
// specific elements.
type Case struct {
	// Namespace refers to the package that the element comes from. This is
	// so different element packages can have elements with the same name
	// while still allowing themes to differentiate between them.
	Namespace string

	// Element refers to the name of the element. This should (generally) be
	// the type name of the element. For example: Button, Input, Container,
	// etc.
	Element string

	// Component specifies the specific part of the element that is being
	// referred to. This parameter is entirely optional.
	Component string
}
 
// C can be used as shorthand to generate a case struct. The component parameter
// may be left out of this argument list for brevity. Arguments passed after
// component will be ignored.
func C (namespace, element string, component ...string) Case {
	if component == nil { component = []string { "" } }
	return Case {
		Namespace: namespace,
		Element:   element,
		Component: component[0],
	}
}

// Match determines if a case matches the specified parameters. A blank string
// will act as a wildcard.
func (c Case) Match (namespace, element, component string) bool {
	if namespace == "" { namespace = c.Namespace }
	if element   == "" { element   = c.Element   }
	if component == "" { component = c.Component }

	return  namespace == c.Namespace &&
		element   == c.Element   &&
		component == c.Component
}

// State lists parameters which can change the appearance of some patterns and
// colors. For example, passing a State with Selected set to true may result in
// a pattern that has a colored border within it.
type State struct {
	// On should be set to true if the element that is using this pattern is
	// in some sort of selected or "on" state, such as if a checkbox is
	// checked, a file is selected, or a switch is toggled on. This is only
	// necessary if the element in question is capable of being toggled or
	// selected.
	On bool

	// Focused should be set to true if the element that is using this
	// pattern is currently focused.
	Focused bool

	// Pressed should be set to true if the element that is using this
	// pattern is being pressed down by the mouse. This is only necessary if
	// the element in question processes mouse button events.
	Pressed bool

	// Disabled should be set to true if the element that is using this
	// pattern is locked and cannot be interacted with. Disabled variations
	// of patterns are typically flattened and greyed-out.
	Disabled bool
}
