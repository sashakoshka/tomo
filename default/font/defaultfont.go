package font

import "golang.org/x/image/font/basicfont"

var FaceRegular = basicfont.Face7x13

// TODO: make bold, italic, and bold italic masks by processing the Face7x13
// mask.

var FaceBold = &basicfont.Face {
	Advance: 7,
	Width:   6,
	Height:  13,
	Ascent:  11,
	Descent: 2,
	Mask:    FaceRegular.Mask, // TODO
	Ranges: []basicfont.Range {
		{ '\u0020', '\u007f', 0  },
		{ '\ufffd', '\ufffe', 95 },
	},
}

var FaceItalic = &basicfont.Face {
	Advance: 7,
	Width:   6,
	Height:  13,
	Ascent:  11,
	Descent: 2,
	Mask:    FaceRegular.Mask, // TODO
	Ranges: []basicfont.Range {
		{ '\u0020', '\u007f', 0  },
		{ '\ufffd', '\ufffe', 95 },
	},
}

var FaceBoldItalic = &basicfont.Face {
	Advance: 7,
	Width:   6,
	Height:  13,
	Ascent:  11,
	Descent: 2,
	Mask:    FaceRegular.Mask, // TODO
	Ranges: []basicfont.Range {
		{ '\u0020', '\u007f', 0  },
		{ '\ufffd', '\ufffe', 95 },
	},
}
