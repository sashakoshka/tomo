package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"

var progressBarCase = tomo.C("tomo", "progressBar")

// ProgressBar displays a visual indication of how far along a task is.
type ProgressBar struct {
	entity tomo.Entity
	progress float64
}

// NewProgressBar creates a new progress bar displaying the given progress
// level.
func NewProgressBar (progress float64) (element *ProgressBar) {
	if progress < 0 { progress = 0 }
	if progress > 1 { progress = 1 }
	element = &ProgressBar { progress: progress }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.updateMinimumSize()
	return
}

// Entity returns this element's entity.
func (element *ProgressBar) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *ProgressBar) Draw (destination artist.Canvas) {
	bounds := element.entity.Bounds()

	pattern := element.entity.Theme().Pattern(tomo.PatternSunken, tomo.State { }, progressBarCase)
	padding := element.entity.Theme().Padding(tomo.PatternSunken, progressBarCase)
	pattern.Draw(destination, bounds)
	bounds = padding.Apply(bounds)
	meterBounds := image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Min.X + int(float64(bounds.Dx()) * element.progress),
		bounds.Max.Y)
	mercury := element.entity.Theme().Pattern(tomo.PatternMercury, tomo.State { }, progressBarCase)
	mercury.Draw(destination, meterBounds)
}

// SetProgress sets the progress level of the bar.
func (element *ProgressBar) SetProgress (progress float64) {
	if progress < 0 { progress = 0 }
	if progress > 1 { progress = 1 }
	if progress == element.progress { return }
	element.progress = progress
	element.entity.Invalidate()
}

func (element *ProgressBar) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *ProgressBar) updateMinimumSize() {
	padding      := element.entity.Theme().Padding(tomo.PatternSunken, progressBarCase)
	innerPadding := element.entity.Theme().Padding(tomo.PatternMercury, progressBarCase)
	element.entity.SetMinimumSize (
		padding.Horizontal() + innerPadding.Horizontal(),
		padding.Vertical()   + innerPadding.Vertical())
}
