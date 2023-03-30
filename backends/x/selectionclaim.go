package x

import "io"
import "github.com/jezek/xgb"
import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xprop"
import "github.com/jezek/xgbutil/xevent"
import "git.tebibyte.media/sashakoshka/tomo/data"

type selectionClaim struct {
	window *window
	data data.Data
	name xproto.Atom
}

func (window *window) claimSelection (name xproto.Atom, data data.Data) *selectionClaim {
	// Follow:
	// https://tronche.com/gui/x/icccm/sec-2.html#s-2.1
	
	// A client wishing to acquire ownership of a particular selection
	// should call SetSelectionOwner. The client should set the specified
	// selection to the atom that represents the selection, set the
	// specified owner to some window that the client created, and set the
	// specified time to some time between the current last-change time of
	// the selection concerned and the current server time. This time value
	// usually will be obtained from the timestamp of the event that
	// triggers the acquisition of the selection. Clients should not set the
	// time value to CurrentTime, because if they do so, they have no way of
	// finding when they gained ownership of the selection. Clients must use
	// a window they created so that requestors can route events to the
	// owner of the selection.
	err := xproto.SetSelectionOwnerChecked (
		window.backend.connection.Conn(),
		window.xWindow.Id, name, 0).Check() // FIXME: should not be zero
	if err != nil { return nil }

	ownerReply, err := xproto.GetSelectionOwner (
		window.backend.connection.Conn(), name).Reply()
	if err != nil { return nil }
	if ownerReply.Owner != window.xWindow.Id { return nil}

	return &selectionClaim {
		window: window,
		data: data,
		name: name,
	}
}

func (window *window) refuseSelectionRequest (request xevent.SelectionRequestEvent) {
	// ... refuse the SelectionRequest by sending the requestor window a
	// SelectionNotify event with the property set to None (by means of a
	// SendEvent request with an empty event mask).
	event := xproto.SelectionNotifyEvent {
		Requestor: request.Requestor,
		Selection: request.Selection,
		Target:    request.Target,
		Property:  0,
	}.Bytes()
	xproto.SendEvent (
		window.backend.connection.Conn(),
		false, request.Requestor, 0, string(event))
}

func (window *window) fulfillSelectionRequest (
	data []byte,
	format byte,
	request xevent.SelectionRequestEvent,
) {
	die := func () { window.refuseSelectionRequest(request) }
	
	// If the specified property is not None, the owner should place the
	// data resulting from converting the selection into the specified
	// property on the requestor window and should set the property's type
	// to some appropriate value, which need not be the same as the
	// specified target.
	err := xproto.ChangePropertyChecked (
		window.backend.connection.Conn(),
		xproto.PropModeReplace, request.Requestor,
		request.Property,
		request.Target, format,
		uint32(len(data) / (int(format) / 8)), data).Check()
	if err != nil { die() }

	// If the property is successfully stored, the owner should acknowledge
	// the successful conversion by sending the requestor window a
	// SelectionNotify event (by means of a SendEvent request with an empty
	// mask).
	event := xproto.SelectionNotifyEvent {
		Requestor: request.Requestor,
		Selection: request.Selection,
		Target:    request.Target,
		Property:  request.Property,
	}.Bytes()
	xproto.SendEvent (
		window.backend.connection.Conn(),
		false, request.Requestor, 0, string(event))
}

func (claim *selectionClaim) handleSelectionRequest (
	connection *xgbutil.XUtil,
	event xevent.SelectionRequestEvent,
) {
	// Follow:
	// https://tronche.com/gui/x/icccm/sec-2.html#s-2.2

	die := func () { claim.window.refuseSelectionRequest(event) }

	// When a requestor wants the value of a selection, the owner receives a
	// SelectionRequest event. The specified owner and selection will be the
	// values that were specified in the SetSelectionOwner request. The
	// owner should compare the timestamp with the period it has owned the
	// selection and, if the time is outside, refuse the SelectionRequest.
	if event.Selection != claim.name { die(); return }

	// If the specified property is None , the requestor is an obsolete
	// client. Owners are encouraged to support these clients by using the
	// specified target atom as the property name to be used for the reply.
	if event.Property == 0 {
		event.Property = event.Target
	}

	// Otherwise, the owner should use the target to decide the form into
	// which the selection should be converted. Some targets may be defined
	// such that requestors can pass parameters along with the request. The
	// owner will find these parameters in the property named in the
	// selection request. The type, format, and contents of this property
	// are dependent upon the definition of the target. If the target is not
	// defined to have parameters, the owner should ignore the property if
	// it is present. If the selection cannot be converted into a form based
	// on the target (and parameters, if any), the owner should refuse the
	// SelectionRequest as previously described.
	targetName, err := xprop.AtomName (
		claim.window.backend.connection, event.Target)
	if err != nil { die(); return }
	
	switch targetName {
	case "TARGETS":
		targetNames := []string { }
		for mime := range claim.data {
			targetNames = append(targetNames, mimeToTargets(mime)...)
		}
		data := make([]byte, len(targetNames) * 4)
		for index, name := range targetNames {
			atom, err := xprop.Atm(claim.window.backend.connection, name)
			if err != nil { die(); return }
			xgb.Put32(data[(index) * 4:], uint32(atom))
		}
		claim.window.fulfillSelectionRequest(data, 32, event)

	default:
		mime, confidence := targetToMime(targetName)
		if confidence == confidenceNone { die(); return }
		reader, ok := claim.data[mime]
		if !ok { die(); return }
		reader.Seek(0, io.SeekStart)
		data, err := io.ReadAll(reader)
		if err != nil { die() }
		claim.window.fulfillSelectionRequest(data, 8, event)
	}
}
