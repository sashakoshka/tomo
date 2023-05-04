package fun

import "time"
import "math"
import "image"
import "image/color"
import "tomo"
import "art"
import "art/shapes"

var clockCase = tomo.C("tomo", "clock")

// AnalogClock can display the time of day in an analog format.
type AnalogClock struct {
	entity tomo.Entity
	time   time.Time
}

// NewAnalogClock creates a new analog clock that displays the specified time.
func NewAnalogClock (newTime time.Time) (element *AnalogClock) {
	element = &AnalogClock { }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.entity.SetMinimumSize(64, 64)
	return
}

// Entity returns this element's entity.
func (element *AnalogClock) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *AnalogClock) Draw (destination art.Canvas) {
	bounds := element.entity.Bounds()

	state   := tomo.State { }
	pattern := element.entity.Theme().Pattern(tomo.PatternSunken, state, clockCase)
	padding := element.entity.Theme().Padding(tomo.PatternSunken, clockCase)
	pattern.Draw(destination, bounds)

	bounds = padding.Apply(bounds)

	foreground := element.entity.Theme().Color(tomo.ColorForeground, state, clockCase)
	accent     := element.entity.Theme().Color(tomo.ColorAccent, state, clockCase)

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

func (element *AnalogClock) HandleThemeChange () {
	element.entity.Invalidate()
}

func (element *AnalogClock) radialLine (
	destination art.Canvas,
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
