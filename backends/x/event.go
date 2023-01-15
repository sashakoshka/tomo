package x

import "git.tebibyte.media/sashakoshka/tomo"

import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xevent"

type scrollSum struct {
	x, y int
}

func (sum *scrollSum) add (button xproto.Button) {
	switch button {
	case 4:
		sum.y --
	case 5:
		sum.y ++
	case 6:
		sum.x --
	case 7:
		sum.x ++
	}
}

func (window *Window) handleConfigureNotify (
	connection *xgbutil.XUtil,
	event xevent.ConfigureNotifyEvent,
) {
	if window.child == nil { return }

	configureEvent := *event.ConfigureNotifyEvent
	
	newWidth  := int(configureEvent.Width)
	newHeight := int(configureEvent.Height)
	sizeChanged :=
		window.metrics.width  != newWidth ||
		window.metrics.height != newHeight
	window.metrics.width  = newWidth
	window.metrics.height = newHeight

	if sizeChanged {
		configureEvent = window.compressConfigureNotify(configureEvent)
		window.metrics.width  = int(configureEvent.Width)
		window.metrics.height = int(configureEvent.Height)
		window.reallocateCanvas()
		window.resizeChildToFit()
	}
}

func (window *Window) handleKeyPress (
	connection *xgbutil.XUtil,
	event xevent.KeyPressEvent,
) {
	if window.child == nil { return}
	
	keyEvent := *event.KeyPressEvent
	key, numberPad := window.backend.keycodeToKey(keyEvent.Detail, keyEvent.State)
	modifiers := tomo.Modifiers {
		Shift:
			(keyEvent.State & xproto.ModMaskShift)                    > 0 ||
			(keyEvent.State & window.backend.modifierMasks.shiftLock) > 0,
		Control: (keyEvent.State & xproto.ModMaskControl)              > 0,
		Alt:     (keyEvent.State & window.backend.modifierMasks.alt)   > 0,
		Meta:    (keyEvent.State & window.backend.modifierMasks.meta)  > 0,
		Super:   (keyEvent.State & window.backend.modifierMasks.super) > 0,
		Hyper:   (keyEvent.State & window.backend.modifierMasks.hyper) > 0,
		NumberPad: numberPad,
	}

	if key == tomo.KeyTab && modifiers.Alt {
		if _, ok := window.child.(tomo.Selectable); ok {
			direction := tomo.SelectionDirectionForward
			if modifiers.Shift {
				direction = tomo.SelectionDirectionBackward
			}

			window.advanceSelectionInChild(direction)
		}
	} else if child, ok := window.child.(tomo.KeyboardTarget); ok {
		// FIXME: pass correct value for repeated
		child.HandleKeyDown(key, modifiers, false)
	}
}

func (window *Window) advanceSelectionInChild (direction tomo.SelectionDirection) {
	child := window.child.(tomo.Selectable)
	if child.Selected() {
		if !child.HandleSelection(direction) {
			child.HandleDeselection()
		}
	} else {
		child.HandleSelection(tomo.SelectionDirectionNeutral)
	}
}

func (window *Window) handleKeyRelease (
	connection *xgbutil.XUtil,
	event xevent.KeyReleaseEvent,
) {
	if window.child == nil { return }
	
	keyEvent := *event.KeyReleaseEvent
	key, numberPad := window.backend.keycodeToKey(keyEvent.Detail, keyEvent.State)
	modifiers := tomo.Modifiers {
		Shift:
			(keyEvent.State & xproto.ModMaskShift)                    > 0 ||
			(keyEvent.State & window.backend.modifierMasks.shiftLock) > 0,
		Control: (keyEvent.State & xproto.ModMaskControl)              > 0,
		Alt:     (keyEvent.State & window.backend.modifierMasks.alt)   > 0,
		Meta:    (keyEvent.State & window.backend.modifierMasks.meta)  > 0,
		Super:   (keyEvent.State & window.backend.modifierMasks.super) > 0,
		Hyper:   (keyEvent.State & window.backend.modifierMasks.hyper) > 0,
		NumberPad: numberPad,
	}
	
	if child, ok := window.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyUp(key, modifiers)
	}
}

func (window *Window) handleButtonPress (
	connection *xgbutil.XUtil,
	event xevent.ButtonPressEvent,
) {
	if window.child == nil { return }
	
	if child, ok := window.child.(tomo.MouseTarget); ok {
		buttonEvent := *event.ButtonPressEvent
		if buttonEvent.Detail >= 4 && buttonEvent.Detail <= 7 {
			sum := scrollSum { }
			sum.add(buttonEvent.Detail)
			window.compressScrollSum(buttonEvent, &sum)
			child.HandleScroll (
				int(buttonEvent.EventX),
				int(buttonEvent.EventY),
				float64(sum.x), float64(sum.y))
		} else {
			child.HandleMouseDown (
				int(buttonEvent.EventX),
				int(buttonEvent.EventY),
				tomo.Button(buttonEvent.Detail))
		}
	}
	
}

func (window *Window) handleButtonRelease (
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
			tomo.Button(buttonEvent.Detail))
	}
}

func (window *Window) handleMotionNotify (
	connection *xgbutil.XUtil,
	event xevent.MotionNotifyEvent,
) {
	if window.child == nil { return }
	
	if child, ok := window.child.(tomo.MouseTarget); ok {
		motionEvent := window.compressMotionNotify(*event.MotionNotifyEvent)
		child.HandleMouseMove (
			int(motionEvent.EventX),
			int(motionEvent.EventY))
	}
}


func (window *Window) compressConfigureNotify (
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

func (window *Window) compressScrollSum (
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

			sum.add(typedEvent.Detail)
			defer func (index int) {
				xevent.DequeueAt(window.backend.connection, index)
			} (index)
		}
	}

	return
}

func (window *Window) compressMotionNotify (
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

		if firstEvent.Event == typedEvent.Event &&
			typedEvent.Detail >= 4 &&
			typedEvent.Detail <= 7 {

			lastEvent = typedEvent
			defer func (index int) {
				xevent.DequeueAt(window.backend.connection, index)
			} (index)
		}
	}

	return
}
