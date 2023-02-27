package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/artist/patterns"

type Image struct {
	*core.Core
	core core.CoreControl
	buffer canvas.Canvas
}

func NewImage (image image.Image) (element *Image) {
	element = &Image { buffer: canvas.FromImage(image) }
	element.Core, element.core = core.NewCore(element.draw)
	bounds := image.Bounds()
	element.core.SetMinimumSize(bounds.Dx(), bounds.Dy())
	return
}

func (element *Image) draw () {
	(patterns.Texture { Canvas: element.buffer }).
		Draw(element.core, element.Bounds())
}
