package layouts

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Vertical lays its children out vertically. It can contain any number of
// children. When an child is added to the layout, it can either be set to
// contract to its minimum height or expand to fill the remaining space (space
// that is not taken up by other children or padding is divided equally among
// these). Child elements will all have the same width.
type Vertical struct {
	
}
