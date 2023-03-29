package x

import "errors"
import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xprop"
import "github.com/jezek/xgbutil/xevent"
import "git.tebibyte.media/sashakoshka/tomo/data"

type selReqState int; const (
	selReqStateClosed selReqState = iota
	selReqStateAwaitSelectionNotify
)

type selectionRequest struct {
	state       selReqState
	window      *window
	source      xproto.Atom
	destination xproto.Atom
	accept      []data.Mime
	callback func (data.Data, error)
}

// TODO: take in multiple formats and check the TARGETS list against them.

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

	// TODO: account for all types in accept slice
	targetName := request.accept[0].String()
	if request.accept[0] == data.M("text", "plain") {
		targetName = "UTF8_STRING"
	}
	targetAtom, err := xprop.Atm(window.backend.connection, targetName)
	if err != nil { request.die(err); return }

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
	err = xproto.DeletePropertyChecked (
		window.backend.connection.Conn(),
		window.xWindow.Id,
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
		targetAtom,
		request.destination,
		// TODO: *possibly replace this zero with an actual timestamp
		// received from the server. this is non-trivial as we cannot
		// rely on the timestamp of the last received event, because
		// there is a possibility that this method is invoked
		// asynchronously from within tomo.Do().
		0).Check()
	if err != nil { request.die(err); return }

	request.state = selReqStateAwaitSelectionNotify
	return
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

func (request *selectionRequest) handleSelectionNotify (
	connection *xgbutil.XUtil,
	event xevent.SelectionNotifyEvent,
) {
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

	// FIXME: get the mime type from the selection owner's response
	request.finalize(data.Bytes(request.accept[0],reply.Value))
}
