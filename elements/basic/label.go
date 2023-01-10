package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type Label struct {
	*core.Core
	core core.CoreControl
	
	text   string
	drawer artist.TextDrawer
}

func NewLabel (text string) (element *Label) {
	element = &Label {  }
	element.Core, element.core = core.NewCore(element)
	face := theme.FontFaceRegular()
	element.drawer.SetFace(face)
	element.SetText(text)
	metrics := face.Metrics()
	emspace, _ := face.GlyphAdvance('M')
	intEmspace := emspace.Round()
	if intEmspace < 1 { intEmspace = theme.Padding()}
	element.core.SetMinimumSize(intEmspace, metrics.Height.Round())
	return
}

func (element *Label) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		element.drawer.SetMaxWidth (resizeEvent.Width)
		element.drawer.SetMaxHeight(resizeEvent.Height)
		element.draw()
	}
	return
}

func (element *Label) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText(text)
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *Label) Selectable () (selectable bool) {
	return
}

func (element *Label) AdvanceSelection (direction int) (ok bool) {
	return
}

func (element *Label) draw () {
	bounds := element.core.Bounds()

	artist.Rectangle (
		element.core,
		theme.BackgroundImage(),
		nil, 0,
		bounds)

	textBounds := element.drawer.LayoutBounds()

	foreground := theme.ForegroundImage()
	element.drawer.Draw (element.core, foreground, image.Point {
		X: 0 - textBounds.Min.X,
		Y: 0 - textBounds.Min.Y,
	})
}
