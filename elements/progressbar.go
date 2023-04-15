package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// ProgressBar displays a visual indication of how far along a task is.
type ProgressBar struct {
	entity tomo.Entity
	progress float64
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewProgressBar creates a new progress bar displaying the given progress
// level.
func NewProgressBar (progress float64) (element *ProgressBar) {
	element = &ProgressBar { progress: progress }
	element.entity = tomo.NewEntity(element)
	element.theme.Case = tomo.C("tomo", "progressBar")
	element.updateMinimumSize()
	return
}

// Entity returns this element's entity.
func (element *ProgressBar) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *ProgressBar) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()

	pattern := element.theme.Pattern(tomo.PatternSunken, tomo.State { })
	padding := element.theme.Padding(tomo.PatternSunken)
	pattern.Draw(destination, bounds)
	bounds = padding.Apply(bounds)
	meterBounds := image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Min.X + int(float64(bounds.Dx()) * element.progress),
		bounds.Max.Y)
	mercury := element.theme.Pattern(tomo.PatternMercury, tomo.State { })
	mercury.Draw(destination, meterBounds)
}

// SetProgress sets the progress level of the bar.
func (element *ProgressBar) SetProgress (progress float64) {
	if progress == element.progress { return }
	element.progress = progress
	element.entity.Invalidate()
}

// SetTheme sets the element's theme.
func (element *ProgressBar) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *ProgressBar) SetConfig (new tomo.Config) {
	if new == nil || new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *ProgressBar) updateMinimumSize() {
	padding      := element.theme.Padding(tomo.PatternSunken)
	innerPadding := element.theme.Padding(tomo.PatternMercury)
	element.entity.SetMinimumSize (
		padding.Horizontal() + innerPadding.Horizontal(),
		padding.Vertical()   + innerPadding.Vertical())
}
