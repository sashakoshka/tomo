package fun

import "time"
import "math"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// AnalogClock can display the time of day in an analog format.
type AnalogClock struct {
	*core.Core
	core core.CoreControl
	time time.Time
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewAnalogClock creates a new analog clock that displays the specified time.
func NewAnalogClock (newTime time.Time) (element *AnalogClock) {
	element = &AnalogClock { }
	element.theme.Case = tomo.C("tomo", "clock")
	element.Core, element.core = core.NewCore(element, element.draw)
	element.core.SetMinimumSize(64, 64)
	return
}

// SetTime changes the time that the clock displays.
func (element *AnalogClock) SetTime (newTime time.Time) {
	if newTime == element.time { return }
	element.time = newTime
	element.redo()
}

// SetTheme sets the element's theme.
func (element *AnalogClock) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *AnalogClock) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.redo()
}

func (element *AnalogClock) redo () {
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *AnalogClock) draw () {
	bounds := element.Bounds()

	state   := tomo.State { }
	pattern := element.theme.Pattern(tomo.PatternSunken, state)
	padding := element.theme.Padding(tomo.PatternSunken)
	pattern.Draw(element.core, bounds)

	bounds = padding.Apply(bounds)

	foreground := element.theme.Color(tomo.ColorForeground, state)
	accent     := element.theme.Color(tomo.ColorAccent, state)

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

func (element *AnalogClock) radialLine (
	source color.RGBA,
	inner  float64,
	outer  float64,
	radian float64,
) {
	bounds := element.Bounds()
	width  := float64(bounds.Dx()) / 2
	height := float64(bounds.Dy()) / 2
	min := element.Bounds().Min.Add(image.Pt (
		int(math.Cos(radian) * inner * width + width),
		int(math.Sin(radian) * inner * height + height)))
	max := element.Bounds().Min.Add(image.Pt (
		int(math.Cos(radian) * outer * width + width),
		int(math.Sin(radian) * outer * height + height)))
	shapes.ColorLine(element.core, source, 1, min, max)
}
