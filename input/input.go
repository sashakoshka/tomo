package input

import "unicode"

// Key represents a keyboard key.
type Key int

const (
	KeyNone Key = 0
	
	KeyInsert      Key = 1
	KeyMenu        Key = 2
	KeyPrintScreen Key = 3
	KeyPause       Key = 4
	KeyCapsLock    Key = 5
	KeyScrollLock  Key = 6
	KeyNumLock     Key = 7
	KeyBackspace   Key = 8
	KeyTab         Key = 9
	KeyEnter       Key = 10
	KeyEscape      Key = 11
	
	KeyUp       Key = 12
	KeyDown     Key = 13
	KeyLeft     Key = 14
	KeyRight    Key = 15
	KeyPageUp   Key = 16
	KeyPageDown Key = 17
	KeyHome     Key = 18
	KeyEnd      Key = 19
	
	KeyLeftShift    Key = 20
	KeyRightShift   Key = 21
	KeyLeftControl  Key = 22
	KeyRightControl Key = 23
	KeyLeftAlt      Key = 24
	KeyRightAlt     Key = 25
	KeyLeftMeta     Key = 26
	KeyRightMeta    Key = 27	
	KeyLeftSuper    Key = 28
	KeyRightSuper   Key = 29
	KeyLeftHyper    Key = 30
	KeyRightHyper   Key = 31
	
	KeyDelete Key = 127
	
	KeyDead Key = 128
	
	KeyF1  Key = 129
	KeyF2  Key = 130
	KeyF3  Key = 131
	KeyF4  Key = 132
	KeyF5  Key = 133
	KeyF6  Key = 134
	KeyF7  Key = 135
	KeyF8  Key = 136
	KeyF9  Key = 137
	KeyF10 Key = 138
	KeyF11 Key = 139
	KeyF12 Key = 140
)

// Button represents a mouse button.
type Button int

const (
	ButtonNone Button = iota
	
	Button1
	Button2
	Button3
	Button4
	Button5
	Button6
	Button7
	Button8
	Button9
	
	ButtonLeft    Button = Button1
	ButtonMiddle  Button = Button2
	ButtonRight   Button = Button3
	ButtonBack    Button = Button8
	ButtonForward Button = Button9
)

// Printable returns whether or not the key's character could show up on screen.
// If this function returns true, the key can be cast to a rune and used as
// such.
func (key Key) Printable () (printable bool) {
	printable = unicode.IsPrint(rune(key))
	return
}

// Modifiers lists what modifier keys are being pressed. This is used in
// conjunction with a Key code in a Key press event. These should be used
// instead of attempting to track the state of the modifier keys, because there
// is no guarantee that one press event will be coupled with one release event.
type Modifiers struct {
	Shift   bool
	Control bool
	Alt     bool
	Meta    bool
	Super   bool
	Hyper   bool

	// NumberPad does not represent a key, but it behaves like one. If it is
	// set to true, the Key was pressed on the number pad. It is treated
	// as a modifier key because if you don't care whether a key was pressed
	// on the number pad or not, you can just ignore this value.
	NumberPad bool
}

// KeynavDirection represents a keyboard navigation direction.
type KeynavDirection int

const (
	KeynavDirectionNeutral  KeynavDirection =  0
	KeynavDirectionBackward KeynavDirection = -1
	KeynavDirectionForward  KeynavDirection =  1
)

// Canon returns a well-formed direction.
func (direction KeynavDirection) Canon () (canon KeynavDirection) {
	if direction > 0 {
		return KeynavDirectionForward
	} else if direction == 0 {
		return KeynavDirectionNeutral
	} else {
		return KeynavDirectionBackward
	}
}
