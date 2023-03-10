package theme

import "image"
import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// IconSize is a type representing valid icon sizes.
type IconSize int

const (
	IconSizeSmall IconSize = 16
	IconSizeLarge IconSize = 48
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
)

// Color lits a number of cannonical colors, each with its own ID.
type Color int; const (
	// ColorAccent is the accent color of the theme.
	ColorAccent Color = iota

	// ColorForeground is the text/icon color of the theme.
	ColorForeground
)

// Icon lists a number of cannonical icons, each with its own ID.
type Icon int; const (
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
	IconSettings)

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
	IconRefresh
	IconHistory

	IconYes
	IconNo)

const (
	// Status icons
	IconInformation Icon = iota + 0x180
	IconQuestion
	IconWarning
	IconError)

const (
	// Tool icons
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
