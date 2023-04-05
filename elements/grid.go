package elements

import "image"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"
// import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type gridCell struct {
	rune
	tomo.FontStyle
	background tomo.Color
	foreground tomo.Color
	clean bool
}

type gridBuffer struct {
	cells  []gridCell
	stride int
}

// Grid is an array of monospaced character cells. Each one has a foreground and
// background color. It satisfies io.Writer and can be fed text with ANSI escape
// codes.
type Grid struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	
	cells      []gridCell
	stride     int
	cellWidth  int
	cellHeight int

	face font.Face
	config config.Wrapped
	theme  theme.Wrapped
	
	onResize func ()
}

func NewGrid () (element *Grid) {
	element = &Grid { }
	element.theme.Case = tomo.C("tomo", "grid")
	element.Core, element.core = core.NewCore(element, element.drawAndPush)
	element.updateFont()
	element.updateMinimumSize()
	return
}

func (element *Grid) HandleMouseDown (x, y int, button input.Button) {
	
}

func (element *Grid) HandleMouseUp (x, y int, button input.Button) {
	
}

func (element *Grid) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	
}

func (element *Grid) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	
}

// SetTheme sets the element's theme.
func (element *Grid) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.updateFont()
	element.updateMinimumSize()
	element.drawAndPush()
}

// SetConfig sets the element's configuration.
func (element *Grid) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.drawAndPush()
}

func (element *Grid) alloc () bool {
	bounds := element.Bounds()
	width  := bounds.Dx() / element.cellWidth
	height := bounds.Dy() / element.cellHeight
	unchanged :=
		width  == element.stride &&
		height == len(element.cells) / element.stride
	if unchanged { return false }

	// TODO: attempt to wrap text
	element.stride = width
	element.cells  = make([]gridCell, width * height)
	for index := range element.cells {
		element.cells[index].background = tomo.ColorBackground
		element.cells[index].foreground = tomo.ColorForeground
	}
	return true
}

func (element *Grid) updateFont () {
	element.face = element.theme.FontFace (
		tomo.FontStyleMonospace,
		tomo.FontSizeNormal)
	emSpace, _ := element.face.GlyphAdvance('M')
	metrics    := element.face.Metrics()
	element.cellWidth  = emSpace.Round()
	element.cellHeight = metrics.Height.Round()
}

func (element *Grid) updateMinimumSize () {
	element.core.SetMinimumSize(element.cellWidth, element.cellHeight)
}

func (element *Grid) state () tomo.State {
	return tomo.State {
		
	}
}

func (element *Grid) drawAndPush () {
	if element.core.HasImage () {
		element.core.DamageRegion(element.draw(true))
	}
}

func (element *Grid) draw (force bool) image.Rectangle {

}
