package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// ProgressBar displays a visual indication of how far along a task is.
type ProgressBar struct {
	*core.Core
	core core.CoreControl
	progress float64
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewProgressBar creates a new progress bar displaying the given progress
// level.
func NewProgressBar (progress float64) (element *ProgressBar) {
	element = &ProgressBar { progress: progress }
	element.theme.Case = theme.C("basic", "progressBar")
	element.Core, element.core = core.NewCore(element.draw)
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

// SetTheme sets the element's theme.
func (element *ProgressBar) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.updateMinimumSize()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *ProgressBar) SetConfig (new config.Config) {
	if new == nil || new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.redo()
}

func (element (ProgressBar)) updateMinimumSize() {
	element.core.SetMinimumSize (
		element.config.Padding() * 2,
		element.config.Padding() * 2)
}

func (element *ProgressBar) redo () {
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *ProgressBar) draw () {
	bounds := element.Bounds()

	pattern := element.theme.Pattern (
		theme.PatternSunken,
		theme.PatternState { })
	inset := element.theme.Inset(theme.PatternSunken)
	artist.FillRectangle(element, pattern, bounds)
	bounds = inset.Apply(bounds)
	meterBounds := image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Min.X + int(float64(bounds.Dx()) * element.progress),
		bounds.Max.Y)
	accent := element.theme.Pattern (
		theme.PatternAccent,
		theme.PatternState { })
	artist.FillRectangle(element, accent, meterBounds)
}
