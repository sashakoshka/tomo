package x

import "git.tebibyte.media/sashakoshka/tomo/data"

type selectionClaim struct {
	data data.Data
	scheduledDelete bool
}

func (window *window) newSelectionClaim (data data.Data) *selectionClaim {
	return &selectionClaim{
		data: data,
	}
}

func (claim *selectionClaim) idle () bool {
	// TODO
}

func (claim *selectionClaim) handleSelectionRequest (
	// TODO
) {
	// TODO
}
