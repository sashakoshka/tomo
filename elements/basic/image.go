package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type Image struct {
	*core.Core
	core core.CoreControl
	buffer artist.Pattern
}

func NewImage (image image.Image) (element *Image) {
	element = &Image { buffer: artist.NewTexture(image) }
	element.Core, element.core = core.NewCore(element.draw)
	bounds := image.Bounds()
	element.core.SetMinimumSize(bounds.Dx(), bounds.Dy())
	return
}

func (element *Image) draw () {
	artist.FillRectangle(element.core, element.buffer, element.Bounds())
}
