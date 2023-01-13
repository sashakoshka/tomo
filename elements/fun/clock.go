package fun

import "time"
import "math"
import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type AnalogClock struct {
	*core.Core
	core core.CoreControl
	time time.Time
}

func NewAnalogClock (newTime time.Time) (element *AnalogClock) {
	element = &AnalogClock { }
	element.Core, element.core = core.NewCore(element)
	element.core.SetMinimumSize(64, 64)
	return
}

func (element *AnalogClock) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		element.draw()
	}
}

func (element *AnalogClock) SetTime (newTime time.Time) {
	if newTime == element.time { return }
	element.time = newTime
	if element.core.HasImage() {
		element.draw()
		element.core.PushAll()
	}
}

func (element *AnalogClock) draw () {
	bounds := element.core.Bounds()

	artist.ChiseledRectangle (
		element.core,
		theme.BackgroundProfile(true),
		bounds)

	for hour := 0; hour < 12; hour ++ {
		element.radialLine (
			theme.ForegroundImage(),
			0.8, 0.9, float64(hour) / 6 * math.Pi)
	}

	second := float64(element.time.Second())
	minute := float64(element.time.Minute()) + second / 60
	hour   := float64(element.time.Hour())   + minute / 60

	element.radialLine (
		theme.ForegroundImage(),
		0, 0.5, (hour - 3) / 6 * math.Pi)
	element.radialLine (
		theme.ForegroundImage(),
		0, 0.7, (minute - 15) / 30 * math.Pi)
	element.radialLine (
		theme.AccentImage(),
		0, 0.7, (second - 15) / 30 * math.Pi)
}

func (element *AnalogClock) radialLine (
	source tomo.Image,
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
