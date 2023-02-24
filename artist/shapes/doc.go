// Package shapes provides some basic shape drawing routines.
//
// A word about patterns:
//
// Most drawing routines have a version that samples from other canvases, and a
// version that samples from a solid color. None of these routines can use
// patterns directly, but it is entirely possible to have a pattern draw to an
// off-screen canvas and then draw a shape based on that canvas. As a little
// bonus, you can save the canvas for later so you don't have to render the
// pattern again when you need to redraw the shape.
package shapes
