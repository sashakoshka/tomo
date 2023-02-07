package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// ProgressBar displays a visual indication of how far along a task is.
type ProgressBar struct {
	*core.Core
	core core.CoreControl
	progress float64
}

// NewProgressBar creates a new progress bar displaying the given progress
// level.
func NewProgressBar (progress float64) (element *ProgressBar) {
	element = &ProgressBar { progress: progress }
	element.Core, element.core = core.NewCore (
		element.draw,
		element.redo,
		element.redo,
		theme.C("basic", "progressBar"))
	return
}

func (element *ProgressBar) redo () {
	element.core.SetMinimumSize (
		element.core.Config().Padding() * 2,
		element.core.Config().Padding() * 2)
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

// SetProgress sets the progress level of the bar.
func (element *ProgressBar) SetProgress (progress float64) {
	if progress == element.progress { return }
	element.progress = progress
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *ProgressBar) draw () {
	bounds := element.Bounds()

	pattern := element.core.Pattern (
		theme.PatternSunken,
		theme.PatternState { })
	inset := element.core.Inset(theme.PatternSunken)
	artist.FillRectangle(element, pattern, bounds)
	bounds = inset.Apply(bounds)
	meterBounds := image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Min.X + int(float64(bounds.Dx()) * element.progress),
		bounds.Max.Y)
	accent := element.core.Pattern (
		theme.PatternSunken,
		theme.PatternState { })
	artist.FillRectangle(element, accent, meterBounds)
}
