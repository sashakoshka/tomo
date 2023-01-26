package basic

import "fmt"
import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// List is an element that contains several objects that a user can select.
type List struct {
	*core.Core
	core core.CoreControl
	
	enabled bool
	selected bool
	pressed bool
	
	contentHeight int
	forcedMinimumWidth  int
	forcedMinimumHeight int
	
	selectedEntry int
	scroll int
	entries []ListEntry
	
	onSelectionRequest func () (granted bool)
	onSelectionMotionRequest func (tomo.SelectionDirection) (granted bool)
	onScrollBoundsChange func ()
	onNoEntrySelected func ()
}

// NewList creates a new list element with the specified entries.
func NewList (entries ...ListEntry) (element *List) {
	element = &List { enabled: true, selectedEntry: -1 }
	element.Core, element.core = core.NewCore(element)
	
	element.entries = make([]ListEntry, len(entries))
	for index, entry := range entries {
		element.entries[index] = entry
	}
	
	element.updateMinimumSize()
	return
}

// Resize changes the element's size.
func (element *List) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	
	for index, entry := range element.entries {
		element.entries[index] = element.resizeEntryToFit(entry)
	}

	element.draw()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// Collapse forces a minimum width and height upon the list. If a zero value is
// given for a dimension, its minimum will be determined by the list's content.
// If the list's height goes beyond the forced size, it will need to be accessed
// via scrolling. If an entry's width goes beyond the forced size, its text will
// be truncated so that it fits.
func (element *List) Collapse (width, height int) {
	element.forcedMinimumWidth  = width
	element.forcedMinimumHeight = height
	element.updateMinimumSize()
}

func (element *List) HandleMouseDown (x, y int, button tomo.Button) {
	if !element.enabled  { return }
	if !element.selected { element.Select() }
	if button != tomo.ButtonLeft { return }
	element.pressed = true
	if element.selectUnderMouse(x, y) && element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *List) HandleMouseUp (x, y int, button tomo.Button) {
	if button != tomo.ButtonLeft { return }
	element.pressed = false
}

func (element *List) HandleMouseMove (x, y int) {
	if element.pressed {
		if element.selectUnderMouse(x, y) && element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *List) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *List) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if !element.enabled { return }

	altered := false
	switch key {
	case tomo.KeyLeft, tomo.KeyUp:
		altered = element.changeSelectionBy(-1)
		
	case tomo.KeyRight, tomo.KeyDown:
		altered = element.changeSelectionBy(1)

	case tomo.KeyEscape:
		altered = element.selectEntry(-1)
	}
	
	if altered && element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *List) HandleKeyUp(key tomo.Key, modifiers tomo.Modifiers) { }

func (element *List) Selected () (selected bool) {
	return element.selected
}

func (element *List) Select () {
	if !element.enabled { return }
	if element.onSelectionRequest != nil {
		element.onSelectionRequest()
	}
}

func (element *List) HandleSelection (
	direction tomo.SelectionDirection,
) (
	accepted bool,
) {
	direction = direction.Canon()
	if !element.enabled { return false }
	if element.selected && direction != tomo.SelectionDirectionNeutral {
		return false
	}
	
	element.selected = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	return true
}

func (element *List) HandleDeselection () {
	element.selected = false
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *List) OnSelectionRequest (callback func () (granted bool)) {
	element.onSelectionRequest = callback
}

func (element *List) OnSelectionMotionRequest (
	callback func (direction tomo.SelectionDirection) (granted bool),
) {
	element.onSelectionMotionRequest = callback
}

// ScrollContentBounds returns the full content size of the element.
func (element *List) ScrollContentBounds () (bounds image.Rectangle) {
	return image.Rect (
		0, 0,
		1, element.contentHeight)
}

// ScrollViewportBounds returns the size and position of the element's viewport
// relative to ScrollBounds.
func (element *List) ScrollViewportBounds () (bounds image.Rectangle) {
	return image.Rect (
		0, element.scroll,
		0, element.scroll + element.scrollViewportHeight())
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *List) ScrollTo (position image.Point) {
	element.scroll = position.Y
	if element.scroll < 0 {
		element.scroll = 0
	} else if element.scroll > element.maxScrollHeight() {
		element.scroll = element.maxScrollHeight()
	}
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// ScrollAxes returns the supported axes for scrolling.
func (element *List) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *List) scrollViewportHeight () (height int) {
	return element.Bounds().Dy() - theme.Padding()
}

func (element *List) maxScrollHeight () (height int) {
	height =
		element.contentHeight -
		element.scrollViewportHeight()
	if height < 0 { height = 0 }
	return
}

func (element *List) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

// SetEnabled sets whether this list can be interacted with or not.
func (element *List) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

// OnNoEntrySelected sets a function to be called when the user chooses to
// deselect the current selected entry by clicking on empty space within the
// list or by pressing the escape key.
func (element *List) OnNoEntrySelected (callback func ()) {
	element.onNoEntrySelected = callback
}

// CountEntries returns the amount of entries in the list.
func (element *List) CountEntries () (count int) {
	return len(element.entries)
}

// Append adds an entry to the end of the list.
func (element *List) Append (entry ListEntry) {
	// append
	entry.Collapse(element.forcedMinimumWidth)
	element.entries = append(element.entries, entry)

	// recalculate, redraw, notify
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// EntryAt returns the entry at the specified index. If the index is out of
// bounds, it panics.
func (element *List) EntryAt (index int) (entry ListEntry) {
	if index < 0 || index >= len(element.entries) {
		panic(fmt.Sprint("basic.List.EntryAt index out of range: ", index))
	}
	return element.entries[index]
}

// Insert inserts an entry into the list at the speified index. If the index is
// out of bounds, it is constrained either to zero or len(entries).
func (element *List) Insert (index int, entry ListEntry) {
	if index < 0 { index = 0 }
	if index > len(element.entries) { index = len(element.entries) }

	// insert
	element.entries = append (
		element.entries[:index + 1],
		element.entries[index:]...)
	entry.Collapse(element.forcedMinimumWidth)
	element.entries[index] = entry

	// recalculate, redraw, notify
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// Remove removes the entry at the specified index. If the index is out of
// bounds, it panics.
func (element *List) Remove (index int) {
	if index < 0 || index >= len(element.entries) {
		panic(fmt.Sprint("basic.List.Remove index out of range: ", index))
	}

	// delete
	element.entries = append (
		element.entries[:index],
		element.entries[index + 1:]...)

	// recalculate, redraw, notify
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// Replace replaces the entry at the specified index with another. If the index
// is out of bounds, it panics.
func (element *List) Replace (index int, entry ListEntry) {
	if index < 0 || index >= len(element.entries) {
		panic(fmt.Sprint("basic.List.Replace index out of range: ", index))
	}

	// replace
	entry.Collapse(element.forcedMinimumWidth)
	element.entries[index] = entry

	// redraw
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *List) selectUnderMouse (x, y int) (updated bool) {
	bounds := element.Bounds()
	mousePoint := image.Pt(x, y)
	dot := image.Pt (
		bounds.Min.X + theme.Padding() / 2,
		bounds.Min.Y - element.scroll + theme.Padding() / 2)
	
	newlySelectedEntryIndex := -1
	for index, entry := range element.entries {
		entryPosition := dot
		dot.Y += entry.Bounds().Dy()
		if entryPosition.Y > bounds.Max.Y { break }
		if mousePoint.In(entry.Bounds().Add(entryPosition)) {
			newlySelectedEntryIndex = index
			break
		}
	}

	return element.selectEntry(newlySelectedEntryIndex)
}

func (element *List) selectEntry (index int) (updated bool) {
	if element.selectedEntry == index { return false }
	element.selectedEntry = index
	if element.selectedEntry < 0 {
		if element.onNoEntrySelected != nil {
			element.onNoEntrySelected()
		}
	} else {
		element.entries[element.selectedEntry].RunSelect()
	}
	return true
}

func (element *List) changeSelectionBy (delta int) (updated bool) {
	newIndex := element.selectedEntry + delta
	if newIndex < 0 { newIndex = len(element.entries) - 1 }
	if newIndex >= len(element.entries) { newIndex = 0 }
	return element.selectEntry(newIndex)
}

func (element *List) resizeEntryToFit (entry ListEntry) (resized ListEntry) {
	entry.Collapse(element.forcedMinimumWidth)
	return entry
}

func (element *List) updateMinimumSize () {
	element.contentHeight = 0
	for _, entry := range element.entries {
		element.contentHeight += entry.Bounds().Dy()
	}

	minimumWidth  := element.forcedMinimumWidth
	minimumHeight := element.forcedMinimumHeight

	if minimumWidth == 0 {
		minimumWidth = theme.Padding()
		for _, entry := range element.entries {
			entryWidth := entry.Bounds().Dx()
			if entryWidth > minimumWidth {
				minimumWidth = entryWidth
			}
		}
	}

	if minimumHeight == 0 {
		minimumHeight = element.contentHeight + theme.Padding()
	}

	element.core.SetMinimumSize(minimumWidth, minimumHeight)
}

func (element *List) draw () {
	bounds := element.Bounds()

	artist.FillRectangle (
		element,
		theme.ListPattern(element.selected),
		bounds)

	dot := image.Point {
		bounds.Min.X,
		bounds.Min.Y - element.scroll + theme.Padding() / 2,
	}
	for index, entry := range element.entries {
		entryPosition := dot
		dot.Y += entry.Bounds().Dy()
		if dot.Y < bounds.Min.Y { continue }
		if entryPosition.Y > bounds.Max.Y { break }

		if element.selectedEntry == index {
			artist.FillRectangle (
				element,
				theme.ListEntryPattern(true),
				entry.Bounds().Add(entryPosition))
		}
		entry.Draw (
			element, entryPosition,
			element.selectedEntry == index && element.selected)
	}
}
