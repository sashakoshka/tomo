package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

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
	element.theme.Case = tomo.C("tomo", "progressBar")
	element.Core, element.core = core.NewCore(element, element.draw)
	element.updateMinimumSize()
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
func (element *ProgressBar) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.updateMinimumSize()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *ProgressBar) SetConfig (new tomo.Config) {
	if new == nil || new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.redo()
}

func (element *ProgressBar) updateMinimumSize() {
	padding      := element.theme.Padding(tomo.PatternSunken)
	innerPadding := element.theme.Padding(tomo.PatternMercury)
	element.core.SetMinimumSize (
		padding.Horizontal() + innerPadding.Horizontal(),
		padding.Vertical()   + innerPadding.Vertical())
}

func (element *ProgressBar) redo () {
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *ProgressBar) draw () {
	bounds := element.Bounds()

	pattern := element.theme.Pattern(tomo.PatternSunken, tomo.State { })
	padding := element.theme.Padding(tomo.PatternSunken)
	pattern.Draw(element.core, bounds)
	bounds = padding.Apply(bounds)
	meterBounds := image.Rect (
		bounds.Min.X, bounds.Min.Y,
		bounds.Min.X + int(float64(bounds.Dx()) * element.progress),
		bounds.Max.Y)
	mercury := element.theme.Pattern(tomo.PatternMercury, tomo.State { })
	mercury.Draw(element.core, meterBounds)
}
