package tomo

import "tomo/artist"

// Element represents a basic on-screen object. Extended element interfaces are
// defined in the ability module.
type Element interface {
	// Draw causes the element to draw to the specified canvas. The bounds
	// of this canvas specify the area that is actually drawn to, while the
	// Entity bounds specify the actual area of the element.
	Draw (artist.Canvas)

	// Entity returns this element's entity.
	Entity () Entity
}
