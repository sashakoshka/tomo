package shatter

import "image"

// Shatter takes in a bounding rectangle, and several rectangles to be
// subtracted from it. It returns a slice of rectangles that tile together to
// make up the difference between them. This is intended to be used for figuring
// out which areas of a container element's background are covered by other
// elements so it doesn't waste CPU cycles drawing to those areas.
func Shatter (
	glass image.Rectangle,
	rocks ...image.Rectangle,
) (
	tiles []image.Rectangle,
) {
	// in this function, the metaphor of throwing several rocks at a sheet
	// of glass is used to illustrate the concept.

	tiles = []image.Rectangle { glass }
	for _, rock := range rocks {
		
		// check each tile to see if the rock has collided with it
		tileLen := len(tiles)
		for tileIndex := 0; tileIndex < tileLen; tileIndex ++ {
			tile := tiles[tileIndex]
			if !rock.Overlaps(tile) { continue }
			newTiles, n := shatterOnce(tile, rock)
			if n > 0 {
				// the tile was shattered into one or more sub
				// tiles
				tiles[tileIndex] = newTiles[0]
				tiles = append(tiles, newTiles[1:n]...)
			} else {
				// the tile was entirely obscured by the rock
				// and must be wholly removed
				tiles = remove(tiles, tileIndex)
				tileIndex --
				tileLen --
			}
		}
	}
	return
}

func shatterOnce (glass, rock image.Rectangle) (tiles [4]image.Rectangle, n int) {
	rock = rock.Intersect(glass)

	// |'''''''''''|
	// |           |
	// |###|'''|   |
	// |###|___|   |
	// |           |
	// |___________|
	if rock.Min.X > glass.Min.X { tiles[n] = image.Rect (
		glass.Min.X, rock.Min.Y,
		rock.Min.X,  rock.Max.Y,
	); n ++ }
	
	// |'''''''''''|
	// |           |
	// |   |'''|###|
	// |   |___|###|
	// |           |
	// |___________|
	if rock.Max.X < glass.Max.X { tiles[n] = image.Rect (
		rock.Max.X,  rock.Min.Y,
		glass.Max.X, rock.Max.Y,
	); n ++ }
	
	// |###########|
	// |###########|
	// |   |'''|   |
	// |   |___|   |
	// |           |
	// |___________|
	if rock.Min.Y > glass.Min.Y { tiles[n] = image.Rect (
		glass.Min.X, glass.Min.Y,
		glass.Max.X, rock.Min.Y,
	); n ++ }
	
	// |'''''''''''|
	// |           |
	// |   |'''|   |
	// |   |___|   |
	// |###########|
	// |###########|
	if rock.Max.Y < glass.Max.Y { tiles[n] = image.Rect (
		glass.Min.X, rock.Max.Y,
		glass.Max.X, glass.Max.Y,
	); n ++ }
	return
}

func remove[ELEMENT any] (slice []ELEMENT, s int) []ELEMENT {
    return append(slice[:s], slice[s + 1:]...)
}
