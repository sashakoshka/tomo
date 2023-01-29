package fun

import "time"
import "math"
import "image"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// AnalogClock can display the time of day in an analog format.
type AnalogClock struct {
	*core.Core
	core core.CoreControl
	time time.Time
}

// NewAnalogClock creates a new analog clock that displays the specified time.
func NewAnalogClock (newTime time.Time) (element *AnalogClock) {
	element = &AnalogClock { }
	element.Core, element.core = core.NewCore(element)
	element.core.SetMinimumSize(64, 64)
	return
}

// Resize changes the size of the clock.
func (element *AnalogClock) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.draw()
}

// SetTime changes the time that the clock displays.
func (element *AnalogClock) SetTime (newTime time.Time) {
	if newTime == element.time { return }
	element.time = newTime
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *AnalogClock) draw () {
	bounds := element.core.Bounds()

	pattern, inset := theme.SunkenPattern(theme.PatternState { })
	artist.FillRectangle(element, pattern, bounds)

	bounds = inset.Apply(bounds)

	foreground, _ := theme.ForegroundPattern(theme.PatternState { })
	accent,     _ := theme.AccentPattern(theme.PatternState { })

	for hour := 0; hour < 12; hour ++ {
		element.radialLine (
			foreground,
			0.8, 0.9, float64(hour) / 6 * math.Pi)
	}

	second := float64(element.time.Second())
	minute := float64(element.time.Minute()) + second / 60
	hour   := float64(element.time.Hour())   + minute / 60

	element.radialLine(foreground, 0, 0.5, (hour   - 3)  / 6  * math.Pi)
	element.radialLine(foreground, 0, 0.7, (minute - 15) / 30 * math.Pi)
	element.radialLine(accent,     0, 0.7, (second - 15) / 30 * math.Pi)
}

// FlexibleHeightFor constrains the clock's minimum size to a 1:1 aspect ratio.
func (element *AnalogClock) FlexibleHeightFor (width int) (height int) {
	return width
}

// OnFlexibleHeightChange sets a function to be calle dwhen the parameters
// affecting the clock's flexible height change.
func (element *AnalogClock) OnFlexibleHeightChange (func ()) { }

func (element *AnalogClock) radialLine (
	source artist.Pattern,
	inner  float64,
	outer  float64,
	radian float64,
) {
	bounds := element.core.Bounds()
	width  := float64(bounds.Dx()) / 2
	height := float64(bounds.Dy()) / 2
	min := image.Pt (
		int(math.Cos(radian) * inner * width + width),
		int(math.Sin(radian) * inner * height + height))
	max := image.Pt (
		int(math.Cos(radian) * outer * width + width),
		int(math.Sin(radian) * outer * height + height))
	// println(min.String(), max.String())
	artist.Line(element.core, source, 1, min, max)
}
