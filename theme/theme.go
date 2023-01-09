package theme

import "image"
import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"

// none of these colors are final! TODO: generate these values from a theme
// file at startup.

var foregroundImage  = artist.NewUniform(color.Gray16 { 0x0000})
var disabledForegroundImage = artist.NewUniform(color.Gray16 { 0x5555})
var accentImage      = artist.NewUniform(color.RGBA { 0x3E, 0x81, 0x69, 0xFF})
var highlightImage   = artist.NewUniform(color.Gray16 { 0xEEEE })
var shadowImage      = artist.NewUniform(color.Gray16 { 0x3333 })
var weakShadeImage   = artist.NewUniform(color.Gray16 { 0x7777 })
var strokeImage      = artist.NewUniform(color.Gray16 { 0x0000 })
var weakStrokeImage  = artist.NewUniform(color.Gray16 { 0x3333 })
var insetShadowImage = artist.NewUniform(color.Gray16 { 0x7777 })

var backgroundImage = artist.NewUniform(color.Gray16 { 0xAAAA})
var backgroundProfile = artist.ShadingProfile {
	Highlight:     highlightImage,
	Shadow:        shadowImage,
	Stroke:        artist.NewUniform(color.Gray16 { 0x0000 }),
	Fill:          backgroundImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var engravedBackgroundProfile = backgroundProfile.Engraved()

var raisedImage = artist.NewUniform(color.RGBA { 0x8D, 0x98, 0x94, 0xFF})
var raisedProfile = artist.ShadingProfile {
	Highlight:     highlightImage,
	Shadow:        shadowImage,
	Stroke:        strokeImage,
	Fill:          raisedImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var selectedRaisedProfile = artist.ShadingProfile {
	Highlight:     highlightImage,
	Shadow:        shadowImage,
	Stroke:        accentImage,
	Fill:          raisedImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var engravedRaisedProfile = artist.ShadingProfile {
	Highlight:     weakShadeImage,
	Shadow:        raisedImage,
	Stroke:        strokeImage,
	Fill:          raisedImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var selectedEngravedRaisedProfile = artist.ShadingProfile {
	Highlight:     insetShadowImage,
	Shadow:        raisedImage,
	Stroke:        accentImage,
	Fill:          raisedImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var disabledRaisedProfile = artist.ShadingProfile {
	Highlight:     weakShadeImage,
	Shadow:        weakShadeImage,
	Stroke:        weakStrokeImage,
	Fill:          backgroundImage,
	StrokeWeight:  1,
	ShadingWeight: 0,
}

var inputImage = artist.NewUniform(color.Gray16 { 0xFFFF })
var inputProfile = artist.ShadingProfile {
	Highlight:     insetShadowImage,
	Shadow:        inputImage,
	Stroke:        strokeImage,
	Fill:          inputImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var selectedInputProfile = artist.ShadingProfile {
	Highlight:     insetShadowImage,
	Shadow:        inputImage,
	Stroke:        accentImage,
	Fill:          inputImage,
	StrokeWeight:  1,
	ShadingWeight: 1,
}
var disabledInputProfile = artist.ShadingProfile {
	Highlight:     weakShadeImage,
	Shadow:        backgroundImage,
	Stroke:        accentImage,
	Fill:          backgroundImage,
	StrokeWeight:  1,
	ShadingWeight: 0,
}

// BackgroundProfile returns the shading profile to be used for backgrounds.
func BackgroundProfile (engraved bool) artist.ShadingProfile {
	if engraved {
		return engravedBackgroundProfile
	} else {
		return backgroundProfile
	}
}

// RaisedProfile returns the shading profile to be used for raised objects such
// as buttons.
func RaisedProfile (
	engraved bool,
	enabled  bool,
	selected bool,
) (
	artist.ShadingProfile,
) {
	if enabled {
		if engraved {
			if selected {
				return selectedEngravedRaisedProfile
			} else {
				return engravedRaisedProfile
			}
		} else {
			if selected {
				return selectedRaisedProfile
			} else {
				return raisedProfile
			}
		}
	} else {
		return disabledRaisedProfile
	}
}

// InputProfile returns the shading profile to be used for input fields.
func InputProfile (enabled bool, selected bool) artist.ShadingProfile {
	if enabled {
		if selected {
			return selectedInputProfile
		} else {
			return inputProfile
		}
	} else {
		return disabledInputProfile
	}
}

// BackgroundImage returns the texture/color used for the fill of
// BackgroundProfile.
func BackgroundImage () tomo.Image {
	return backgroundImage
}

// RaisedImage returns the texture/color used for the fill of RaisedProfile.
func RaisedImage () tomo.Image {
	return raisedImage
}

// InputImage returns the texture/color used for the fill of InputProfile.
func InputImage () tomo.Image {
	return inputImage
}

// ForegroundImage returns the texture/color text and monochromatic icons should
// be drawn with.
func ForegroundImage () tomo.Image {
	return foregroundImage
}

// DisabledForegroundImage returns the texture/color text and monochromatic
// icons should be drawn with if they are disabled.
func DisabledForegroundImage () tomo.Image {
	return disabledForegroundImage
}

// AccentImage returns the accent texture/color.
func AccentImage () tomo.Image {
	return accentImage
}

// TODO: load fonts from an actual source instead of using basicfont

// FontFaceRegular returns the font face to be used for normal text.
func FontFaceRegular () font.Face {
	return defaultfont.FaceRegular
}

// FontFaceBold returns the font face to be used for bolded text.
func FontFaceBold () font.Face {
	return defaultfont.FaceBold
}

// FontFaceItalic returns the font face to be used for italicized text.
func FontFaceItalic () font.Face {
	return defaultfont.FaceItalic
}

// FontFaceBoldItalic returns the font face to be used for text that is both
// bolded and italicized.
func FontFaceBoldItalic () font.Face {
	return defaultfont.FaceBoldItalic
}

// Padding returns how spaced out things should be on the screen. Generally,
// text should be offset from its container on all sides by this amount.
func Padding () int {
	return 8
}

// SinkOffsetVector specifies a vector for things such as text to move by when a
// "sinking in" effect is desired, such as a button label during a button press.
func SinkOffsetVector () image.Point {
	return image.Point { 1, 1 }
}
