package artist

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"

// ShadingProfile contains shading information that can be used to draw chiseled
// objects.
type ShadingProfile struct {
	Highlight     tomo.Image
	Shadow        tomo.Image
	Stroke        tomo.Image
	Fill          tomo.Image
	StrokeWeight  int
	ShadingWeight int
}

// Engraved reverses the shadown and highlight colors of the ShadingProfile to
// produce a new ShadingProfile with an engraved appearance.
func (profile ShadingProfile) Engraved () (reversed ShadingProfile) {
	reversed = profile
	reversed.Highlight = profile.Shadow
	reversed.Shadow    = profile.Highlight
	return
}

// ChiseledRectangle draws a rectangle with a chiseled/embossed appearance,
// according to the ShadingProfile passed to it.
func ChiseledRectangle (
	destination tomo.Canvas,
	profile     ShadingProfile,
	bounds      image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	// FIXME: this breaks when the bounds are smaller than the border or
	// shading weight

	stroke    := profile.Stroke
	highlight := profile.Highlight
	shadow    := profile.Shadow
	fill      := profile.Fill
	strokeWeight  := profile.StrokeWeight
	shadingWeight := profile.ShadingWeight

	bounds = bounds.Canon()
	updatedRegion = bounds

	strokeWeightVector  := image.Point { strokeWeight,  strokeWeight  }
	shadingWeightVector := image.Point { shadingWeight, shadingWeight }

	shadingBounds := bounds
	shadingBounds.Min = shadingBounds.Min.Add(strokeWeightVector)
	shadingBounds.Max = shadingBounds.Max.Sub(strokeWeightVector)
	shadingBounds = shadingBounds.Canon()

	fillBounds := shadingBounds
	fillBounds.Min = fillBounds.Min.Add(shadingWeightVector)
	fillBounds.Max = fillBounds.Max.Sub(shadingWeightVector)
	fillBounds = fillBounds.Canon()

	strokeImageMin    := stroke.Bounds().Min
	highlightImageMin := highlight.Bounds().Min
	shadowImageMin    := shadow.Bounds().Min
	fillImageMin      := fill.Bounds().Min

	width  := float64(bounds.Dx())
	height := float64(bounds.Dy())

	yy := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
		xx := 0
		for x := bounds.Min.X; x < bounds.Max.X; x ++ {
			var pixel color.RGBA
			point := image.Point { x, y }
			switch {
			case point.In(fillBounds):
				pixel = fill.RGBAAt (
					xx - strokeWeight - shadingWeight +
					fillImageMin.X,
					yy - strokeWeight - shadingWeight +
					fillImageMin.Y)
					
			case point.In(shadingBounds):
				var highlighted bool
				// FIXME: this doesn't work quite right, the
				// slope of the line is somewhat off.
				bottomCorner :=
					float64(xx) < float64(yy) *
					(width / height)
				if bottomCorner {
					highlighted =
						float64(xx) <
						height - float64(yy)
				} else {
					highlighted =
						width - float64(xx) >
						float64(yy)
				}
			
				if highlighted {
					pixel = highlight.RGBAAt (
						xx - strokeWeight +
						highlightImageMin.X,
						yy - strokeWeight +
						highlightImageMin.Y)
				} else {
					pixel = shadow.RGBAAt (
						xx - strokeWeight +
						shadowImageMin.X,
						yy - strokeWeight +
						shadowImageMin.Y)
				}
				
			default:
				pixel = stroke.RGBAAt (
					xx + strokeImageMin.X,
					yy + strokeImageMin.Y)
			}
			destination.SetRGBA(x, y, pixel)
			xx ++
		}
		yy ++
	}
	
	return
}
