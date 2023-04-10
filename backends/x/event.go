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
		case 4:
			sum.x -= scrollDistance
		case 5:
			sum.x += scrollDistance
		case 6:
			sum.y -= scrollDistance
		case 7:
			sum.y += scrollDistance
		}
	} else {
		switch button {
		case 4:
			sum.y -= scrollDistance
		case 5:
			sum.y += scrollDistance
		case 6:
			sum.x -= scrollDistance
		case 7:
			sum.x += scrollDistance
		}
	}

}

func (window *window) handleExpose (
	connection *xgbutil.XUtil,
	event xevent.ExposeEvent,
) {
	_, region := window.compressExpose(*event.ExposeEvent)
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
			window.redrawChildEntirely()
		}
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
	if window.child == nil { return }
	if window.hasModal     { return }
	
	keyEvent := *event.KeyPressEvent
	key, numberPad := window.backend.keycodeToKey(keyEvent.Detail, keyEvent.State)
	modifiers := window.modifiersFromState(keyEvent.State)
	modifiers.NumberPad = numberPad

	if key == input.KeyTab && modifiers.Alt {
		if child, ok := window.child.(tomo.Focusable); ok {
			direction := input.KeynavDirectionForward
			if modifiers.Shift {
				direction = input.KeynavDirectionBackward
			}

			if !child.HandleFocus(direction) {
				child.HandleUnfocus()
			}
		}
	} else if child, ok := window.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyDown(key, modifiers)
	}
}

func (window *window) handleKeyRelease (
	connection *xgbutil.XUtil,
	event xevent.KeyReleaseEvent,
) {
	if window.child == nil { return }
	
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
	
	if child, ok := window.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyUp(key, modifiers)
	}
}

func (window *window) handleButtonPress (
	connection *xgbutil.XUtil,
	event xevent.ButtonPressEvent,
) {
	if window.child == nil { return }
	if window.hasModal     { return }
	
	buttonEvent := *event.ButtonPressEvent
	if buttonEvent.Detail >= 4 && buttonEvent.Detail <= 7 {
		if child, ok := window.child.(tomo.ScrollTarget); ok {
			sum := scrollSum { }
			sum.add(buttonEvent.Detail, window, buttonEvent.State)
			window.compressScrollSum(buttonEvent, &sum)
			child.HandleScroll (
				int(buttonEvent.EventX),
				int(buttonEvent.EventY),
				float64(sum.x), float64(sum.y))
		}
	} else {
		if child, ok := window.child.(tomo.MouseTarget); ok {
			child.HandleMouseDown (
				int(buttonEvent.EventX),
				int(buttonEvent.EventY),
				input.Button(buttonEvent.Detail))
		}
	}
	
}

func (window *window) handleButtonRelease (
	connection *xgbutil.XUtil,
	event xevent.ButtonReleaseEvent,
) {
	if window.child == nil { return }
	
	if child, ok := window.child.(tomo.MouseTarget); ok {
		buttonEvent := *event.ButtonReleaseEvent
		if buttonEvent.Detail >= 4 && buttonEvent.Detail <= 7 { return }
		child.HandleMouseUp (
			int(buttonEvent.EventX),
			int(buttonEvent.EventY),
			input.Button(buttonEvent.Detail))
	}
}

func (window *window) handleMotionNotify (
	connection *xgbutil.XUtil,
	event xevent.MotionNotifyEvent,
) {
	if window.child == nil { return }
	
	if child, ok := window.child.(tomo.MotionTarget); ok {
		motionEvent := window.compressMotionNotify(*event.MotionNotifyEvent)
		child.HandleMotion (
			int(motionEvent.EventX),
			int(motionEvent.EventY))
	}
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
