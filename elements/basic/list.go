package basicElements

import "fmt"
import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// List is an element that contains several objects that a user can select.
type List struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl

	pressed bool
	
	contentHeight int
	forcedMinimumWidth  int
	forcedMinimumHeight int
	
	selectedEntry int
	scroll int
	entries []ListEntry
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onScrollBoundsChange func ()
	onNoEntrySelected func ()
}

// NewList creates a new list element with the specified entries.
func NewList (entries ...ListEntry) (element *List) {
	element = &List { selectedEntry: -1 }
	element.theme.Case = theme.C("basic", "list")
	element.Core, element.core = core.NewCore(element.handleResize)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore (func () {
		if element.core.HasImage () {
			element.draw()
			element.core.DamageAll()
		}
	})
	
	element.entries = make([]ListEntry, len(entries))
	for index, entry := range entries {
		element.entries[index] = entry
	}
	
	element.updateMinimumSize()
	return
}

func (element *List) handleResize () {
	for index, entry := range element.entries {
		element.entries[index] = element.resizeEntryToFit(entry)
	}

	element.draw()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// SetTheme sets the element's theme.
func (element *List) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	for index, entry := range element.entries {
		entry.SetTheme(element.theme.Theme)
		element.entries[index] = entry
	}
	element.updateMinimumSize()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *List) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	for index, entry := range element.entries {
		entry.SetConfig(element.config)
		element.entries[index] = entry
	}
	element.updateMinimumSize()
	element.redo()
}

func (element *List) redo () {
	for index, entry := range element.entries {
		element.entries[index] = element.resizeEntryToFit(entry)
	}

	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
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
	if
		element.forcedMinimumWidth == width &&
		element.forcedMinimumHeight == height {
		
		return
	}
	
	element.forcedMinimumWidth  = width
	element.forcedMinimumHeight = height
	element.updateMinimumSize()

	for index, entry := range element.entries {
		element.entries[index] = element.resizeEntryToFit(entry)
	}
	
	element.redo()
}

func (element *List) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled()  { return }
	if !element.Focused() { element.Focus() }
	if button != input.ButtonLeft { return }
	element.pressed = true
	if element.selectUnderMouse(x, y) && element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *List) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
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

func (element *List) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }

	altered := false
	switch key {
	case input.KeyLeft, input.KeyUp:
		altered = element.changeSelectionBy(-1)
		
	case input.KeyRight, input.KeyDown:
		altered = element.changeSelectionBy(1)

	case input.KeyEscape:
		altered = element.selectEntry(-1)
	}
	
	if altered && element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *List) HandleKeyUp(key input.Key, modifiers input.Modifiers) { }

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
	inset := element.theme.Inset(theme.PatternSunken)
	return element.Bounds().Dy() - inset[0] - inset[2]
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
	entry = element.resizeEntryToFit(entry)
	entry.SetTheme(element.theme.Theme)
	entry.SetConfig(element.config)
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
	entry = element.resizeEntryToFit(entry)
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
	entry = element.resizeEntryToFit(entry)
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

// Select selects a specific item in the list. If the index is out of bounds,
// no items will be selecected.
func (element *List) Select (index int) {
	if element.selectEntry(index) {
		element.redo()
	}
}

func (element *List) selectUnderMouse (x, y int) (updated bool) {
	inset := element.theme.Inset(theme.PatternSunken)
	bounds := inset.Apply(element.Bounds())
	mousePoint := image.Pt(x, y)
	dot := image.Pt (
		bounds.Min.X,
		bounds.Min.Y - element.scroll)
	
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
	bounds := element.Bounds()
	inset := element.theme.Inset(theme.PatternSunken)
	entry.Resize(inset.Apply(bounds).Dx())
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
		for _, entry := range element.entries {
			entryWidth := entry.MinimumWidth()
			if entryWidth > minimumWidth {
				minimumWidth = entryWidth
			}
		}
	}

	if minimumHeight == 0 {
		minimumHeight = element.contentHeight
	}

	inset := element.theme.Inset(theme.PatternSunken)
	minimumHeight += inset[0] + inset[2]

	element.core.SetMinimumSize(minimumWidth, minimumHeight)
}

func (element *List) draw () {
	bounds := element.Bounds()

	inset := element.theme.Inset(theme.PatternSunken)
	pattern := element.theme.Pattern (theme.PatternSunken, theme.PatternState {
		Disabled: !element.Enabled(),
		Focused: element.Focused(),
	})
	artist.FillRectangle(element, pattern, bounds)

	bounds = inset.Apply(bounds)
	dot := image.Point {
		bounds.Min.X,
		bounds.Min.Y - element.scroll,
	}
	innerCanvas := canvas.Cut(element, bounds)
	for index, entry := range element.entries {
		entryPosition := dot
		dot.Y += entry.Bounds().Dy()
		if dot.Y < bounds.Min.Y { continue }
		if entryPosition.Y > bounds.Max.Y { break }
		entry.Draw (
			innerCanvas, entryPosition,
			element.Focused(), element.selectedEntry == index)
	}
}
