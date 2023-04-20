package fun

import "time"
import "math"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

// AnalogClock can display the time of day in an analog format.
type AnalogClock struct {
	entity tomo.Entity
	time   time.Time
	theme  theme.Wrapped
}

// NewAnalogClock creates a new analog clock that displays the specified time.
func NewAnalogClock (newTime time.Time) (element *AnalogClock) {
	element = &AnalogClock { }
	element.theme.Case = tomo.C("tomo", "clock")
	element.entity = tomo.NewEntity(element)
	element.entity.SetMinimumSize(64, 64)
	return
}

// Entity returns this element's entity.
func (element *AnalogClock) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *AnalogClock) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()

	state   := tomo.State { }
	pattern := element.theme.Pattern(tomo.PatternSunken, state)
	padding := element.theme.Padding(tomo.PatternSunken)
	pattern.Draw(destination, bounds)

	bounds = padding.Apply(bounds)

	foreground := element.theme.Color(tomo.ColorForeground, state)
	accent     := element.theme.Color(tomo.ColorAccent, state)

	for hour := 0; hour < 12; hour ++ {
		element.radialLine (
			destination,
			foreground,
			0.8, 0.9, float64(hour) / 6 * math.Pi)
	}

	second := float64(element.time.Second())
	minute := float64(element.time.Minute()) + second / 60
	hour   := float64(element.time.Hour())   + minute / 60

	element.radialLine(destination, foreground, 0, 0.5, (hour   - 3)  / 6  * math.Pi)
	element.radialLine(destination, foreground, 0, 0.7, (minute - 15) / 30 * math.Pi)
	element.radialLine(destination, accent,     0, 0.7, (second - 15) / 30 * math.Pi)
}

// SetTime changes the time that the clock displays.
func (element *AnalogClock) SetTime (newTime time.Time) {
	if newTime == element.time { return }
	element.time = newTime
	element.entity.Invalidate()
}

// SetTheme sets the element's theme.
func (element *AnalogClock) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.entity.Invalidate()
}

func (element *AnalogClock) radialLine (
	destination canvas.Canvas,
	source color.RGBA,
	inner  float64,
	outer  float64,
	radian float64,
) {
	bounds := element.entity.Bounds()
	width  := float64(bounds.Dx()) / 2
	height := float64(bounds.Dy()) / 2
	min := bounds.Min.Add(image.Pt (
		int(math.Cos(radian) * inner * width + width),
		int(math.Sin(radian) * inner * height + height)))
	max := bounds.Min.Add(image.Pt (
		int(math.Cos(radian) * outer * width + width),
		int(math.Sin(radian) * outer * height + height)))
	shapes.ColorLine(destination, source, 1, min, max)
}
