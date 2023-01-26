package basic

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
	element.Core, element.core = core.NewCore(element)
	element.core.SetMinimumSize(theme.Padding() * 2, theme.Padding() * 2)
	return
}

// Resize resizes the progress bar.
func (element *ProgressBar) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.draw()
	return
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
	bounds := element.core.Bounds()

	artist.FillRectangle (
		element.core,
		theme.SunkenPattern(false),
		bounds)
	meterBounds := image.Rect (
		bounds.Min.X + 2, bounds.Min.Y + 2,
		bounds.Min.X - 1 + int(float64(bounds.Dx()) * element.progress),
		bounds.Dy() - 1)
	artist.FillRectangle (
		element.core,
		theme.AccentPattern(),
		meterBounds)
}
