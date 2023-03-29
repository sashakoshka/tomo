package x

import "errors"
import "strings"
import "github.com/jezek/xgb"
import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xprop"
import "github.com/jezek/xgbutil/xevent"
import "git.tebibyte.media/sashakoshka/tomo/data"

type selReqState int; const (
	selReqStateClosed selReqState = iota
	selReqStateAwaitTargets
	selReqStateAwaitValue
	selReqStateAwaitChunk
)

type selectionRequest struct {
	state       selReqState
	window      *window
	source      xproto.Atom
	destination xproto.Atom
	accept      []data.Mime
	mime        data.Mime
	callback func (data.Data, error)
}

func (window *window) newSelectionRequest (
	source, destination xproto.Atom,
	callback func (data.Data, error),
	accept ...data.Mime,
) (
	request *selectionRequest,
) {
	request = &selectionRequest {
		source:      source,
		destination: destination,
		window:      window,
		accept:      accept,
		callback:    callback,
	}

	targets, err := xprop.Atm(window.backend.connection, "TARGETS")
	if err != nil { request.die(err); return }
	request.convertSelection(targets, selReqStateAwaitTargets)
	return
}

func (request *selectionRequest) convertSelection (
	target xproto.Atom, switchTo selReqState,
) {
	// The requestor should set the property argument to the name of a
	// property that the owner can use to report the value of the selection.
	// Requestors should ensure that the named property does not exist on
	// the window before issuing the ConvertSelection. The exception to this
	// rule is when the requestor intends to pass parameters with the
	// request. Some targets may be defined such that requestors can pass
	// parameters along with the request. If the requestor wishes to provide
	// parameters to a request, they should be placed in the specified
	// property on the requestor window before the requestor issues the
	// ConvertSelection request, and this property should be named in the
	// request.
	err := xproto.DeletePropertyChecked (
		request.window.backend.connection.Conn(),
		request.window.xWindow.Id,
		request.destination).Check()
	if err != nil { request.die(err); return }

	// The selection argument specifies the particular selection involved,
	// and the target argument specifies the required form of the
	// information. For information about the choice of suitable atoms to
	// use, see section 2.6. The requestor should set the requestor argument
	// to a window that it created; the owner will place the reply property
	// there. The requestor should set the time argument to the timestamp on
	// the event that triggered the request for the selection value. Note
	// that clients should not specify CurrentTime*.
	err = xproto.ConvertSelectionChecked (
		request.window.backend.connection.Conn(),
		request.window.xWindow.Id,
		request.source,
		target,
		request.destination,
		// TODO: *possibly replace this zero with an actual timestamp
		// received from the server. this is non-trivial as we cannot
		// rely on the timestamp of the last received event, because
		// there is a possibility that this method is invoked
		// asynchronously from within tomo.Do().
		0).Check()
	if err != nil { request.die(err); return }
	
	request.state = switchTo
}

func (request *selectionRequest) die (err error) {
	request.callback(nil, err)
	request.state = selReqStateClosed
}

func (request *selectionRequest) finalize (data data.Data) {
	request.callback(data, nil)
	request.state = selReqStateClosed
}

func (request *selectionRequest) open () bool {
	return request.state != selReqStateClosed
}

type confidence int; const (
	confidenceNone confidence = iota
	confidencePartial
	confidenceFull
)

func targetToMime (name string) (data.Mime, confidence) {
	// TODO: add stuff like PDFs, etc. reference this table:
	// https://tronche.com/gui/x/icccm/sec-2.html#s-2.6.2
	// perhaps we should also have parameters for mime types so we can
	// return an encoding here for things like STRING?
	switch name {
	case "UTF8_STRING":
		return data.MimePlain, confidenceFull
	case "TEXT":
		return data.MimePlain, confidencePartial
	case "STRING":
		return data.MimePlain, confidencePartial
	default:
		if strings.Count(name, "/") == 1 {
			ty, subtype, _ := strings.Cut(name, "/")
			return data.M(ty, subtype), confidenceFull
		} else {
			return data.Mime { }, confidenceNone
		}
	}
}

func (request *selectionRequest) handleSelectionNotify (
	connection *xgbutil.XUtil,
	event xevent.SelectionNotifyEvent,
) {
	// the only valid states that we can process a SelectionNotify event in
	if request.state != selReqStateAwaitValue && request.state != selReqStateAwaitTargets {
		return
	}
	
	// Follow:
	// https://tronche.com/gui/x/icccm/sec-2.html#s-2.4
	
	// If the property argument is None, the conversion has been refused.
	// This can mean either that there is no owner for the selection, that
	// the owner does not support the conversion implied by the target, or
	// that the server did not have sufficient space to accommodate the
	// data.
	if event.Property == 0 { request.die(nil); return }

	// TODO: handle INCR

	// When using GetProperty to retrieve the value of a selection, the
	// property argument should be set to the corresponding value in the
	// SelectionNotify event. Because the requestor has no way of knowing
	// beforehand what type the selection owner will use, the type argument
	// should be set to AnyPropertyType. Several GetProperty requests may be
	// needed to retrieve all the data in the selection; each should set the
	// long-offset argument to the amount of data received so far, and the
	// size argument to some reasonable buffer size (see section 2.5). If
	// the returned value of bytes-after is zero, the whole property has
	// been transferred.
	reply, err := xproto.GetProperty (
		connection.Conn(), false, request.window.xWindow.Id,
		event.Property, xproto.GetPropertyTypeAny,
		0, (1 << 32) - 1).Reply()
	if err != nil { request.die(err); return }
	if reply.Format == 0 {
		request.die(errors.New("x: missing selection property"))
		return
	}

	// Once all the data in the selection has been retrieved (which may
	// require getting the values of several properties &emdash; see section
	// 2.7), the requestor should delete the property in the SelectionNotify
	// request by using a GetProperty request with the delete argument set
	// to True. As previously discussed, the owner has no way of knowing
	// when the data has been transferred to the requestor unless the
	// property is removed.
	if err != nil { request.die(err); return }
	err = xproto.DeletePropertyChecked (
		request.window.backend.connection.Conn(),
		request.window.xWindow.Id,
		request.destination).Check()
	if err != nil { request.die(err); return }

	switch request.state {
	case selReqStateAwaitValue:
		// we now have the full selection data in the property, so we
		// finalize the request and are done.
		// FIXME: get the type from the property and convert that to the
		// mime value to pass to the application.
		request.finalize(data.Bytes(request.mime, reply.Value))
		
	case selReqStateAwaitTargets:
		// make a list of the atoms we got
		buffer := reply.Value
		atoms  := make([]xproto.Atom, len(buffer) / 4)
		for index := range atoms {
			atoms[index] = xproto.Atom(xgb.Get32(buffer[index * 4:]))
		}

		// choose the best match out of all targets using a confidence
		// system
		confidentMatchFound := false
		var chosenTarget xproto.Atom
		var chosenMime   data.Mime
		for _, atom := range atoms {
			targetName, err := xprop.AtomName (
				request.window.backend.connection, atom)
			if err != nil { request.die(err); return }
			
			mime, confidence := targetToMime(targetName)
			if confidence == confidenceNone { continue }

			// if the accepted types list is nil, just choose this
			// one. however, if we are not 100% confident that this
			// target can be directly converted into a mime type,
			// don't mark it as the final match. we still want the
			// mime type we give to the application to be as
			// accurate as possible.
			if request.accept == nil {
				chosenTarget = atom
				chosenMime   = mime
				if confidence == confidenceFull {
					confidentMatchFound = true
				}
			}

			// run through the accepted types list if it exists,
			// looking for a match. if one is found, then choose
			// this target. however, if we are not 100% confident
			// that this target directly corresponds to the mime
			// type, don't mark it as the final match, because there
			// may be a better target in the list.
			for _, accept := range request.accept {
			if accept == mime {
				chosenTarget = atom
				chosenMime   = mime
				if confidence == confidenceFull {
					confidentMatchFound = true
				}
				break
			}}
			
			if confidentMatchFound { break }
		}

		// if we didn't find a match, finalize the request with an empty
		// data map to inform the application that, although there were
		// no errors, there wasn't a suitable target to choose from.
		if chosenTarget == 0 {
			request.finalize(data.Data { })
			return
		}

		// await the selection value
		request.mime = chosenMime
		request.convertSelection(chosenTarget, selReqStateAwaitValue)
	}
}
