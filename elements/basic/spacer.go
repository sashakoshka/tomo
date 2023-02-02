package basicElements

import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

var spacerCase = theme.C("basic", "spacer")

// Spacer can be used to put space between two elements..
type Spacer struct {
	*core.Core
	core core.CoreControl
	line bool
}

// NewSpacer creates a new spacer. If line is set to true, the spacer will be
// filled with a line color, and if compressed to its minimum width or height,
// will appear as a line.
func NewSpacer (line bool) (element *Spacer) {
	element = &Spacer { line: line }
	element.Core, element.core = core.NewCore(element.draw)
	element.core.SetMinimumSize(1, 1)
	return
}

/// SetLine sets whether or not the spacer will appear as a colored line.
func (element *Spacer) SetLine (line bool) {
	if element.line == line { return }
	element.line = line
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Spacer) draw () {
	bounds := element.Bounds()

	if element.line {
		pattern, _ := theme.ForegroundPattern(theme.PatternState {
			Case: spacerCase,
			Disabled: true,
		})
		artist.FillRectangle(element, pattern, bounds)
	} else {
		pattern, _ := theme.BackgroundPattern(theme.PatternState {
			Case: spacerCase,
			Disabled: true,
		})
		artist.FillRectangle(element, pattern, bounds)
	}
}
