package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

type tableCell struct {
	tomo.Element
	tomo.Pattern
	image.Rectangle
}

// TableContainer is a container that lays its contents out in a table. It can
// be scrolled.
type TableContainer struct {
	*core.Core
	*core.Propagator
	core core.CoreControl
	
	topHeading  bool
	leftHeading bool

	columns  int
	rows     int
	scroll   image.Point
	warping  bool
	grid     [][]tableCell
	children []tomo.Element
	
	contentBounds image.Rectangle
	forcedMinimumWidth  int
	forcedMinimumHeight int

	selectedColumn int
	selectedRow int
	
	config config.Wrapped
	theme  theme.Wrapped

	onSelect func ()
	onScrollBoundsChange func ()
}

func NewTableContainer (
	columns, rows int,
	topHeading, leftHeading bool,
) (
	element *TableContainer,
) {
	element = &TableContainer {
		topHeading:     topHeading,
		leftHeading:    leftHeading,
		selectedColumn: -1,
		selectedRow:    -1,
	}
	
	element.theme.Case = tomo.C("tomo", "tableContainer")
	element.Core, element.core = core.NewCore(element, element.redoAll)
	element.Propagator = core.NewPropagator(element, element.core)
	element.Resize(columns, rows)
	return
}

// Set places an element at the specified column and row. If the element passed
// is nil, whatever element occupies the cell currently is removed.
func (element *TableContainer) Set (column, row int, child tomo.Element) {
	if row    < 0 || row    >= element.rows    { return }
	if column < 0 || column >= element.columns { return }

	childList := element.children
	if child == nil {
		if element.grid[row][column].Element == nil {
			// no-op
			return
		} else {
			// removing the child that is currently in a slow
			element.unhook(element.grid[row][column].Element)
			childList = childList[:len(childList) - 1]
			element.grid[row][column].Element = child
		}
	} else {
		element.hook(child)
		if element.grid[row][column].Element == nil {
			// putting the child in an empty slot
			childList = append(childList, nil)
			element.grid[row][column].Element = child
		} else {
			// replacing the child that is currently in a slow
			element.unhook(element.grid[row][column].Element)
			element.grid[row][column].Element = child
		}
	}

	element.rebuildChildList(childList)
	element.children = childList
	element.redoAll()
}

// Resize changes the amount of columns and rows in the table. If the table is
// resized to be smaller, children in cells that do not exist anymore will be
// removed. The minimum size for a TableContainer is 1x1.
func (element *TableContainer) Resize (columns, rows int) {
	if columns < 1 { columns = 1 }
	if rows    < 1 { rows    = 1 }
	if element.columns == columns && element.rows == rows { return }
	amountRemoved := 0

	// handle rows as a whole
	if rows < element.rows {
		// disown children in bottom rows
		for _,     row   := range element.grid[rows:] {
		for index, child := range row {
		if child.Element != nil {
			element.unhook(child.Element)
			amountRemoved ++
			row[index].Element = nil
		}}}
		// cut grid to size
		element.grid = element.grid[:rows]
	} else {
		// expand grid
		newGrid := make([][]tableCell, rows)
		copy(newGrid, element.grid)
		element.grid = newGrid
	}

	// handle each row individually
	for rowIndex, row := range element.grid {
		if columns < element.columns {
			// disown children in the far right of the row
			for index, child := range row[columns:] {
			if child.Element != nil {
				element.unhook(child.Element)
				amountRemoved ++
				row[index].Element = nil
			}}
			// cut row to size
			element.grid[rowIndex] = row[:columns]
		} else {
			// expand row
			newRow := make([]tableCell, columns)
			copy(newRow, row)
			element.grid[rowIndex] = newRow
		}
	}

	element.columns = columns
	element.rows    = rows

	if amountRemoved > 0 {
		childList := element.children[:len(element.children) - amountRemoved]
		element.rebuildChildList(childList)
		element.children = childList
	}
	element.redoAll()
}

// Selected returns the column and row of the cell that is currently selected.
// If no cell is selected, this method will return (-1, -1).
func (element *TableContainer) Selected () (column, row int) {
	return element.selectedColumn, element.selectedRow
}

// OnSelect sets a function to be called when the user selects a table cell.
func (element *TableContainer) OnSelect (callback func ()) {
	element.onSelect = callback
}

// Warp runs the specified callback, deferring all layout and rendering updates
// until the callback has finished executing. This allows for aplications to
// perform batch gui updates without flickering and stuff.
func (element *TableContainer) Warp (callback func ()) {
	if element.warping {
		callback()
		return
	}

	element.warping = true
	callback()
	element.warping = false
	
	element.redoAll()
}

// Collapse collapses the element's minimum width and height. A value of zero
// for either means that the element's normal value is used.
func (element *TableContainer) Collapse (width, height int) {
	if
		element.forcedMinimumWidth == width &&
		element.forcedMinimumHeight == height {
		
		return
	}
	
	element.forcedMinimumWidth  = width
	element.forcedMinimumHeight = height
	element.updateMinimumSize()
}

// CountChildren returns the amount of children contained within this element.
func (element *TableContainer) CountChildren () (count int) {
	return len(element.children)
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (element *TableContainer) Child (index int) (child tomo.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index]
}

func (element *TableContainer) Window () tomo.Window {
	return element.core.Window()
}

// NotifyMinimumSizeChange notifies the container that the minimum size of a
// child element has changed.
func (element *TableContainer) NotifyMinimumSizeChange (child tomo.Element) {
	element.updateMinimumSize()
	element.redoAll()
}

// DrawBackground draws a portion of the container's background pattern within
// the specified bounds. The container will not push these changes.
func (element *TableContainer) DrawBackground (bounds image.Rectangle) {
	if !bounds.Overlaps(element.core.Bounds()) { return }
	
	for rowIndex,    row   := range element.grid {
	for columnIndex, child := range row {
	if bounds.Overlaps(child.Rectangle) {
		element.theme.Pattern (
			child.Pattern,
			element.state(columnIndex, rowIndex)).
			Draw(canvas.Cut(element.core, bounds), child.Rectangle)
		return
	}}}
}

func (element *TableContainer) HandleMouseDown (x, y int, button input.Button) {
	element.Propagator.HandleMouseDown(x, y, button)
	if button != input.ButtonLeft { return }
	
	for rowIndex,    row   := range element.grid {
	for columnIndex, child := range row {
	if image.Pt(x, y).In(child.Rectangle) {
		selected :=
			rowIndex == element.selectedRow &&
			columnIndex == element.selectedColumn
		if selected { return }
		oldColumn, oldRow := element.selectedColumn, element.selectedRow
		element.selectedColumn, element.selectedRow = columnIndex, rowIndex
		if oldColumn >= 0 && oldRow >= 0 {
			element.core.DamageRegion(element.redoCell(oldColumn, oldRow))
		}
		element.core.DamageRegion(element.redoCell(columnIndex, rowIndex))
		if element.onSelect != nil {
			element.onSelect()
		}
		return
	}}}
	
}

func (element *TableContainer) hook (child tomo.Element) {
	if child0, ok := child.(tomo.Themeable); ok {
		child0.SetTheme(element.theme.Theme)
	}
	if child0, ok := child.(tomo.Configurable); ok {
		child0.SetConfig(element.config.Config)
	}
	child.SetParent(element)
}

func (element *TableContainer) unhook (child tomo.Element) {
	child.SetParent(nil)
	child.DrawTo(nil, image.Rectangle { }, nil)
}

func (element *TableContainer) rebuildChildList (list []tomo.Element) {
	index := 0
	for _, row := range element.grid {
	for _, child := range row {
		if child.Element == nil { continue }
		list[index] = child.Element
		index ++
	}}
}

func (element *TableContainer) state (column, row int) (state tomo.State) {
	if column == element.selectedColumn && row == element.selectedRow {
		state.On = true
	}
	return
}

func (element *TableContainer) redoCell (column, row int) image.Rectangle {
	padding := element.theme.Padding(tomo.PatternTableCell)
	cell := element.grid[row][column]
	pattern := element.theme.Pattern (
		cell.Pattern, element.state(column, row))
		
	if cell.Element != nil {
		// give child canvas portion
		innerCellBounds := padding.Apply(cell.Rectangle)
		artist.DrawShatter (
			element.core, pattern,
			cell.Rectangle, innerCellBounds)
		cell.DrawTo (
			canvas.Cut(element.core, innerCellBounds),
			innerCellBounds,
			element.childDrawCallback)
	} else {
		// draw cell pattern in empty cells
		pattern.Draw(element.core, cell.Rectangle)
	}
	return cell.Rectangle
}

func (element *TableContainer) redoAll () {
	if element.warping || !element.core.HasImage() {
		element.updateMinimumSize()
		return
	}

	// calculate the minimum size of each column and row
	bounds := element.Bounds()
	var minWidth, minHeight float64
	columnWidths := make([]float64, element.columns)
	rowHeights   := make([]float64, element.rows)
	padding := element.theme.Padding(tomo.PatternTableCell)

	for rowIndex,    row   := range element.grid {
	for columnIndex, child := range row {
		width, height := padding.Horizontal(), padding.Vertical()
		
		if child.Element != nil {
			minWidth, minHeight := child.MinimumSize()
			width  += minWidth
			height += minHeight
			fwidth  := float64(width)
			fheight := float64(height)
			if fwidth > columnWidths[columnIndex] {
				columnWidths[columnIndex] = fwidth
			}
			if fheight > rowHeights[rowIndex] {
				rowHeights[rowIndex] = fheight
			}
		}
	}}

	// scale up those minimum sizes to an actual size.
	// FIXME: replace this with a more accurate algorithm
	for _, width  := range columnWidths { minWidth  += width  }
	for _, height := range rowHeights   { minHeight += height }
	widthRatio  := float64(bounds.Dx()) / minWidth
	heightRatio := float64(bounds.Dy()) / minHeight
	for index := range columnWidths {
		columnWidths[index] *= widthRatio
	}
	for index := range rowHeights {
		rowHeights[index] *= heightRatio
	}

	// cut up canvas
	x := float64(bounds.Min.X)
	y := float64(bounds.Min.Y)
	for rowIndex, row := range element.grid {
		for columnIndex, _ := range row {
			width  := columnWidths[columnIndex]
			height := rowHeights[rowIndex]
			cellBounds := image.Rect (
				int(x), int(y),
				int(x + width), int(y + height))
			
			var id tomo.Pattern
			isHeading :=
				rowIndex == 0 && element.topHeading ||
				columnIndex == 0 && element.leftHeading
			if isHeading {
				id = tomo.PatternTableHead
			} else {
				id = tomo.PatternTableCell
			}
			element.grid[rowIndex][columnIndex].Rectangle = cellBounds
			element.grid[rowIndex][columnIndex].Pattern   = id
			
			element.redoCell(columnIndex, rowIndex)
			x += float64(width)
		}
		
		x = float64(bounds.Min.X)
		y += rowHeights[rowIndex]
	}

	element.core.DamageAll()
	
	// update the minimum size of the element
	element.core.SetMinimumSize(int(minWidth), int(minHeight))
}

func (element *TableContainer) updateMinimumSize () {
	columnWidths := make([]int, element.columns)
	rowHeights   := make([]int, element.rows)
	padding := element.theme.Padding(tomo.PatternTableCell)

	for rowIndex,    row   := range element.grid {
	for columnIndex, child := range row {
		width, height := padding.Horizontal(), padding.Vertical()
		
		if child.Element != nil {
			minWidth, minHeight := child.MinimumSize()
			width  += minWidth
			height += minHeight
			if width > columnWidths[columnIndex] {
				columnWidths[columnIndex] = width
			}
			if height > rowHeights[rowIndex] {
				rowHeights[rowIndex] = height
			}
		}
	}}

	var minWidth, minHeight int
	for _, width  := range columnWidths { minWidth  += width  }
	for _, height := range rowHeights   { minHeight += height }

	element.core.SetMinimumSize(minWidth, minHeight)
}

func (element *TableContainer) childDrawCallback (region image.Rectangle) {
	element.core.DamageRegion(region)
}