package x

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"

import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xevent"

type scrollSum struct {
	x, y int
}

const scrollDistance = 16

func (sum *scrollSum) add (button xproto.Button, window *window, state uint16) {
	shift := 
		(state & xproto.ModMaskShift)                    > 0 ||
		(state & window.backend.modifierMasks.shiftLock) > 0
	if shift {
		switch button {
		case 4: sum.x -= scrollDistance
		case 5: sum.x += scrollDistance
		case 6: sum.y -= scrollDistance
		case 7: sum.y += scrollDistance
		}
	} else {
		switch button {
		case 4: sum.y -= scrollDistance
		case 5: sum.y += scrollDistance
		case 6: sum.x -= scrollDistance
		case 7: sum.x += scrollDistance
		}
	}
}

func (window *window) handleExpose (
	connection *xgbutil.XUtil,
	event xevent.ExposeEvent,
) {
	_, region := window.compressExpose(*event.ExposeEvent)
	window.system.afterEvent()
	window.pushRegion(region)
}

func (window *window) updateBounds (x, y int16, width, height uint16) {
	window.metrics.bounds =
		image.Rect(0, 0, int(width), int(height)).
		Add(image.Pt(int(x), int(y)))
}

func (window *window) handleConfigureNotify (
	connection *xgbutil.XUtil,
	event xevent.ConfigureNotifyEvent,
) {
	if window.child == nil { return }

	configureEvent := *event.ConfigureNotifyEvent
	
	newWidth  := int(configureEvent.Width)
	newHeight := int(configureEvent.Height)
	sizeChanged :=
		window.metrics.bounds.Dx() != newWidth ||
		window.metrics.bounds.Dy() != newHeight
	window.updateBounds (
		configureEvent.X, configureEvent.Y,
		configureEvent.Width, configureEvent.Height)

	if sizeChanged {
		configureEvent = window.compressConfigureNotify(configureEvent)
		window.updateBounds (
			configureEvent.X, configureEvent.Y,
			configureEvent.Width, configureEvent.Height)
		window.reallocateCanvas()
		window.resizeChildToFit()

		if !window.exposeEventFollows(configureEvent) {
			window.child.Invalidate()
			window.child.InvalidateLayout()
		}
		
		window.system.afterEvent()
	}
}

func (window *window) exposeEventFollows (event xproto.ConfigureNotifyEvent) (found bool) {	
	nextEvents := xevent.Peek(window.backend.connection)
	if len(nextEvents) > 0 {
		untypedEvent := nextEvents[0]
		if untypedEvent.Err == nil {
			typedEvent, ok :=
				untypedEvent.Event.(xproto.ConfigureNotifyEvent)
			
			if ok && typedEvent.Window == event.Window {
				return true
			}
		}
	}
	return false
}

func (window *window) modifiersFromState (
	state uint16,
) (
	modifiers input.Modifiers,
) {
	return input.Modifiers {
		Shift:
			(state & xproto.ModMaskShift)                    > 0 ||
			(state & window.backend.modifierMasks.shiftLock) > 0,
		Control: (state & xproto.ModMaskControl)              > 0,
		Alt:     (state & window.backend.modifierMasks.alt)   > 0,
		Meta:    (state & window.backend.modifierMasks.meta)  > 0,
		Super:   (state & window.backend.modifierMasks.super) > 0,
		Hyper:   (state & window.backend.modifierMasks.hyper) > 0,
	}
}

func (window *window) handleKeyPress (
	connection *xgbutil.XUtil,
	event xevent.KeyPressEvent,
) {
	if window.hasModal { return }
	
	keyEvent := *event.KeyPressEvent
	key, numberPad := window.backend.keycodeToKey(keyEvent.Detail, keyEvent.State)
	modifiers := window.modifiersFromState(keyEvent.State)
	modifiers.NumberPad = numberPad

	if key == input.KeyTab && modifiers.Alt {
		if modifiers.Shift {
			window.system.focusPrevious()
		} else {
			window.system.focusNext()
		}
	} else if key == input.KeyEscape && window.shy {
		window.Close()
	} else if window.focused != nil {
		focused, ok := window.focused.element.(tomo.KeyboardTarget)
		if ok { focused.HandleKeyDown(key, modifiers) }
	}
	
	window.system.afterEvent()
}

func (window *window) handleKeyRelease (
	connection *xgbutil.XUtil,
	event xevent.KeyReleaseEvent,
) {
	keyEvent := *event.KeyReleaseEvent

	// do not process this event if it was generated from a key repeat
	nextEvents := xevent.Peek(window.backend.connection)
	if len(nextEvents) > 0 {
		untypedEvent := nextEvents[0]
		if untypedEvent.Err == nil {
			typedEvent, ok :=
				untypedEvent.Event.(xproto.KeyPressEvent)
			
			if ok && typedEvent.Detail == keyEvent.Detail &&
				typedEvent.Event == keyEvent.Event &&
				typedEvent.State == keyEvent.State {

				return
			}
		}
	}
	
	key, numberPad := window.backend.keycodeToKey(keyEvent.Detail, keyEvent.State)
	modifiers := window.modifiersFromState(keyEvent.State)
	modifiers.NumberPad = numberPad

	if window.focused != nil {
		focused, ok := window.focused.element.(tomo.KeyboardTarget)
		if ok { focused.HandleKeyUp(key, modifiers) }
		
		window.system.afterEvent()
	}
}

func (window *window) handleButtonPress (
	connection *xgbutil.XUtil,
	event xevent.ButtonPressEvent,
) {
	if window.hasModal { return }
	
	buttonEvent  := *event.ButtonPressEvent
	point        := image.Pt(int(buttonEvent.EventX), int(buttonEvent.EventY))
	insideWindow := point.In(window.canvas.Bounds())
	scrolling    := buttonEvent.Detail >= 4 && buttonEvent.Detail <= 7

	underneath := window.system.childAt(point)
	
	if !insideWindow && window.shy && !scrolling {
		window.Close()
	} else if scrolling {
		if child, ok := underneath.element.(tomo.ScrollTarget); ok {
			sum := scrollSum { }
			sum.add(buttonEvent.Detail, window, buttonEvent.State)
			window.compressScrollSum(buttonEvent, &sum)
			child.HandleScroll (
				point.X, point.Y,
				float64(sum.x), float64(sum.y))
		}
	} else {
		if child, ok := underneath.element.(tomo.MouseTarget); ok {
			window.system.drags[buttonEvent.Detail] = child
			child.HandleMouseDown (
				point.X, point.Y,
				input.Button(buttonEvent.Detail))
		}
	}
	
	window.system.afterEvent()
}

func (window *window) handleButtonRelease (
	connection *xgbutil.XUtil,
	event xevent.ButtonReleaseEvent,
) {
	buttonEvent := *event.ButtonReleaseEvent
	if buttonEvent.Detail >= 4 && buttonEvent.Detail <= 7 { return }
	child := window.system.drags[buttonEvent.Detail]
	if child != nil {
		child.HandleMouseUp (
			int(buttonEvent.EventX),
			int(buttonEvent.EventY),
			input.Button(buttonEvent.Detail))
	}
	
	window.system.afterEvent()
}

func (window *window) handleMotionNotify (
	connection *xgbutil.XUtil,
	event xevent.MotionNotifyEvent,
) {
	motionEvent := window.compressMotionNotify(*event.MotionNotifyEvent)
	x := int(motionEvent.EventX)
	y :=int(motionEvent.EventY)

	handled := false
	for _, child := range window.system.drags {
		if child, ok := child.(tomo.MotionTarget); ok {
			child.HandleMotion(x, y)
			handled = true
		}
	}

	if !handled {
		child := window.system.childAt(image.Pt(x, y))
		if child, ok := child.element.(tomo.MotionTarget); ok {
			child.HandleMotion(x, y)
		}
	}
	
	window.system.afterEvent()
}

func (window *window) handleSelectionNotify (
	connection *xgbutil.XUtil,
	event xevent.SelectionNotifyEvent,
) {
	if window.selectionRequest == nil { return }
	window.selectionRequest.handleSelectionNotify(connection, event)
	if !window.selectionRequest.open() { window.selectionRequest = nil }
}

func (window *window) handlePropertyNotify (
	connection *xgbutil.XUtil,
	event xevent.PropertyNotifyEvent,
) {
	if window.selectionRequest == nil { return }
	window.selectionRequest.handlePropertyNotify(connection, event)
	if !window.selectionRequest.open() { window.selectionRequest = nil }
}

func (window *window) handleSelectionClear (
	connection *xgbutil.XUtil,
	event xevent.SelectionClearEvent,
) {
	window.selectionClaim = nil
}

func (window *window) handleSelectionRequest (
	connection *xgbutil.XUtil,
	event xevent.SelectionRequestEvent,
) {
	if window.selectionClaim == nil { return }
	window.selectionClaim.handleSelectionRequest(connection, event)
}

func (window *window) compressExpose (
	firstEvent xproto.ExposeEvent,
) (
	lastEvent xproto.ExposeEvent,
	region image.Rectangle,
) {
	region = image.Rect (
		int(firstEvent.X), int(firstEvent.Y),
		int(firstEvent.X + firstEvent.Width),
		int(firstEvent.Y + firstEvent.Height))
	
	window.backend.connection.Sync()
	xevent.Read(window.backend.connection, false)
	lastEvent = firstEvent
	
	for index, untypedEvent := range xevent.Peek(window.backend.connection) {
		if untypedEvent.Err != nil { continue }
		
		typedEvent, ok := untypedEvent.Event.(xproto.ExposeEvent)
		if !ok { continue }

		if firstEvent.Window == typedEvent.Window {
			region = region.Union (image.Rect (
				int(typedEvent.X), int(typedEvent.Y),
				int(typedEvent.X + typedEvent.Width),
				int(typedEvent.Y + typedEvent.Height)))
		
			lastEvent = typedEvent
			defer func (index int) {
				xevent.DequeueAt(window.backend.connection, index)
			} (index)
		}
	}

	return
}

func (window *window) compressConfigureNotify (
	firstEvent xproto.ConfigureNotifyEvent,
) (
	lastEvent xproto.ConfigureNotifyEvent,
) {
	window.backend.connection.Sync()
	xevent.Read(window.backend.connection, false)
	lastEvent = firstEvent
	
	for index, untypedEvent := range xevent.Peek(window.backend.connection) {
		if untypedEvent.Err != nil { continue }
		
		typedEvent, ok := untypedEvent.Event.(xproto.ConfigureNotifyEvent)
		if !ok { continue }
		
		if firstEvent.Event == typedEvent.Event &&
			firstEvent.Window == typedEvent.Window {

			lastEvent = typedEvent
			defer func (index int) {
				xevent.DequeueAt(window.backend.connection, index)
			} (index)
		}
	}

	return
}

func (window *window) compressScrollSum (
	firstEvent xproto.ButtonPressEvent,
	sum *scrollSum,
) {
	window.backend.connection.Sync()
	xevent.Read(window.backend.connection, false)
	
	for index, untypedEvent := range xevent.Peek(window.backend.connection) {
		if untypedEvent.Err != nil { continue }
		
		typedEvent, ok := untypedEvent.Event.(xproto.ButtonPressEvent)
		if !ok { continue }

		if firstEvent.Event == typedEvent.Event &&
			typedEvent.Detail >= 4 &&
			typedEvent.Detail <= 7 {

			sum.add(typedEvent.Detail, window, typedEvent.State)
			defer func (index int) {
				xevent.DequeueAt(window.backend.connection, index)
			} (index)
		}
	}

	return
}

func (window *window) compressMotionNotify (
	firstEvent xproto.MotionNotifyEvent,
) (
	lastEvent xproto.MotionNotifyEvent,
) {
	window.backend.connection.Sync()
	xevent.Read(window.backend.connection, false)
	lastEvent = firstEvent
	
	for index, untypedEvent := range xevent.Peek(window.backend.connection) {
		if untypedEvent.Err != nil { continue }
		
		typedEvent, ok := untypedEvent.Event.(xproto.MotionNotifyEvent)
		if !ok { continue }

		if firstEvent.Event == typedEvent.Event {
			lastEvent = typedEvent
			defer func (index int) {
				xevent.DequeueAt(window.backend.connection, index)
			} (index)
		}
	}

	return
}
