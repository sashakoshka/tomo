package theme

import "git.tebibyte.media/sashakoshka/tomo/artist"

func ScrollGutterPattern (horizontal bool, enabled bool) (artist.Pattern) {
	if enabled {
		return sunkenPattern
	} else {
		return disabledButtonPattern
	}
}

func ScrollBarPattern (horizontal bool, enabled bool) (artist.Pattern) {
	if enabled {
		return buttonPattern
	} else {
		return disabledButtonPattern
	}
}
