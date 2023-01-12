package x

import "unicode"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/keybind"
import "git.tebibyte.media/sashakoshka/tomo"

// when making changes to this file, look at keysymdef.h and
// https://tronche.com/gui/x/xlib/input/keyboard-encoding.html

var buttonCodeTable = map[xproto.Keysym] tomo.Key {
	0xFFFFFF: tomo.KeyNone,

	0xFF63: tomo.KeyInsert,
	0xFF67: tomo.KeyMenu,
	0xFF61: tomo.KeyPrintScreen,
	0xFF6B: tomo.KeyPause,
	0xFFE5: tomo.KeyCapsLock,
	0xFF14: tomo.KeyScrollLock,
	0xFF7F: tomo.KeyNumLock,
	0xFF08: tomo.KeyBackspace,
	0xFF09: tomo.KeyTab,
	0xFE20: tomo.KeyTab,
	0xFF0D: tomo.KeyEnter,
	0xFF1B: tomo.KeyEscape,
	
	0xFF52: tomo.KeyUp,
	0xFF54: tomo.KeyDown,
	0xFF51: tomo.KeyLeft,
	0xFF53: tomo.KeyRight,
	0xFF55: tomo.KeyPageUp,
	0xFF56: tomo.KeyPageDown,
	0xFF50: tomo.KeyHome,
	0xFF57: tomo.KeyEnd,
	
	0xFFE1: tomo.KeyLeftShift,
	0xFFE2: tomo.KeyRightShift,
	0xFFE3: tomo.KeyLeftControl,
	0xFFE4: tomo.KeyRightControl,

	0xFFE7: tomo.KeyLeftMeta,
	0xFFE8: tomo.KeyRightMeta,
	0xFFE9: tomo.KeyLeftAlt,
	0xFFEA: tomo.KeyRightAlt,
	0xFFEB: tomo.KeyLeftSuper,
	0xFFEC: tomo.KeyRightSuper,
	0xFFED: tomo.KeyLeftHyper,
	0xFFEE: tomo.KeyRightHyper,
	
	0xFFFF: tomo.KeyDelete,
	
	0xFFBE: tomo.KeyF1,
	0xFFBF: tomo.KeyF2,
	0xFFC0: tomo.KeyF3,
	0xFFC1: tomo.KeyF4,
	0xFFC2: tomo.KeyF5,
	0xFFC3: tomo.KeyF6,
	0xFFC4: tomo.KeyF7,
	0xFFC5: tomo.KeyF8,
	0xFFC6: tomo.KeyF9,
	0xFFC7: tomo.KeyF10,
	0xFFC8: tomo.KeyF11,
	0xFFC9: tomo.KeyF12,

	// TODO: send this whenever a compose key, dead key, etc is pressed,
	// and then send the resulting character while witholding the key
	// presses that were used to compose it. As far as the program is
	// concerned, a magical key with the final character was pressed and the
	// KeyDead key is just so that the program might provide some visual
	// feedback to the user while input is being waited for.
	0xFF20: tomo.KeyDead,
}

var keypadCodeTable = map[xproto.Keysym] tomo.Key {
	0xff80: tomo.Key(' '),
	0xff89: tomo.KeyTab,
	0xff8d: tomo.KeyEnter,
	0xff91: tomo.KeyF1,
	0xff92: tomo.KeyF2,
	0xff93: tomo.KeyF3,
	0xff94: tomo.KeyF4,
	0xff95: tomo.KeyHome,
	0xff96: tomo.KeyLeft,
	0xff97: tomo.KeyUp,
	0xff98: tomo.KeyRight,
	0xff99: tomo.KeyDown,
	0xff9a: tomo.KeyPageUp,
	0xff9b: tomo.KeyPageDown,
	0xff9c: tomo.KeyEnd,
	0xff9d: tomo.KeyHome,
	0xff9e: tomo.KeyInsert,
	0xff9f: tomo.KeyDelete,
	0xffbd: tomo.Key('='),
	0xffaa: tomo.Key('*'),
	0xffab: tomo.Key('+'),
	0xffac: tomo.Key(','),
	0xffad: tomo.Key('-'),
	0xffae: tomo.Key('.'),
	0xffaf: tomo.Key('/'),

	0xffb0: tomo.Key('0'),
	0xffb1: tomo.Key('1'),
	0xffb2: tomo.Key('2'),
	0xffb3: tomo.Key('3'),
	0xffb4: tomo.Key('4'),
	0xffb5: tomo.Key('5'),
	0xffb6: tomo.Key('6'),
	0xffb7: tomo.Key('7'),
	0xffb8: tomo.Key('8'),
	0xffb9: tomo.Key('9'),
}

// initializeKeymapInformation grabs keyboard mapping information from the X
// server.
func (backend *Backend) initializeKeymapInformation () {
	keybind.Initialize(backend.connection)
	backend.modifierMasks.capsLock   = backend.keysymToMask(0xFFE5)
	backend.modifierMasks.shiftLock  = backend.keysymToMask(0xFFE6)
	backend.modifierMasks.numLock    = backend.keysymToMask(0xFF7F)
	backend.modifierMasks.modeSwitch = backend.keysymToMask(0xFF7E)
	
	backend.modifierMasks.hyper = backend.keysymToMask(0xffed)
	backend.modifierMasks.super = backend.keysymToMask(0xffeb)
	backend.modifierMasks.meta  = backend.keysymToMask(0xffe7)
	backend.modifierMasks.alt   = backend.keysymToMask(0xffe9)
}

// keysymToKeycode converts an X keysym to an X keycode, instead of the other
// way around.
func (backend *Backend) keysymToKeycode (
	symbol xproto.Keysym,
) (
	code xproto.Keycode,
) {
	mapping := keybind.KeyMapGet(backend.connection)
	
	for index, testSymbol := range mapping.Keysyms {
		if testSymbol == symbol {
			code = xproto.Keycode (
				index /
				int(mapping.KeysymsPerKeycode) +
				int(backend.connection.Setup().MinKeycode))
			break
		}
	}
	
	return
}

// keysymToMask returns the X modmask for a given modifier key.
func (backend *Backend) keysymToMask (
	symbol xproto.Keysym,
) (
	mask uint16,
) {
	mask = keybind.ModGet (
		backend.connection,
		backend.keysymToKeycode(symbol))
	
	return
}

// keycodeToButton converts an X keycode to a tomo keycode. It implements a more
// fleshed out version of some of the logic found in xgbutil/keybind/encoding.go
// to get a full keycode to keysym conversion, but eliminates redundant work by
// going straight to a tomo keycode.
func (backend *Backend) keycodeToKey (
	keycode xproto.Keycode,
	state   uint16,
) (
	button    tomo.Key,
	numberPad bool,
) {
	// PARAGRAPH 3
	//
	// A list of KeySyms is associated with each KeyCode. The list is
	// intended to convey the set of symbols on the corresponding key. If
	// the list (ignoring trailing NoSymbol entries) is a single KeySym
	// ``K'', then the list is treated as if it were the list ``K NoSymbol
	// K NoSymbol''. If the list (ignoring trailing NoSymbol entries) is a
	// pair of KeySyms ``K1 K2'', then the list is treated as if it were the
	// list ``K1 K2 K1 K2''. If the list (ignoring trailing NoSymbol
	// entries) is a triple of KeySyms ``K1 K2 K3'', then the list is
	// treated as if it were the list ``K1 K2 K3 NoSymbol''. When an
	// explicit ``void'' element is desired in the list, the value
	// VoidSymbol can be used.
	symbol1 := keybind.KeysymGet(backend.connection, keycode, 0)
	symbol2 := keybind.KeysymGet(backend.connection, keycode, 1)
	symbol3 := keybind.KeysymGet(backend.connection, keycode, 2)
	symbol4 := keybind.KeysymGet(backend.connection, keycode, 3)
	switch {
	case symbol2 == 0 && symbol3 == 0 && symbol4 == 0:
		symbol3 = symbol1
	case symbol3 == 0 && symbol4 == 0:
		symbol3 = symbol1
		symbol4 = symbol2
	case symbol4 == 0:
		symbol4 = 0
	}
	symbol1Rune := keysymToRune(symbol1)
	symbol2Rune := keysymToRune(symbol2)
	symbol3Rune := keysymToRune(symbol3)
	symbol4Rune := keysymToRune(symbol4)

	// PARAGRAPH 4
	//
	// The first four elements of the list are split into two groups of
	// KeySyms. Group 1 contains the first and second KeySyms; Group 2
	// contains the third and fourth KeySyms. Within each group, if the
	// second element of the group is NoSymbol , then the group should be
	// treated as if the second element were the same as the first element,
	// except when the first element is an alphabetic KeySym ``K'' for which
	// both lowercase and uppercase forms are defined. In that case, the
	// group should be treated as if the first element were the lowercase
	// form of ``K'' and the second element were the uppercase form of
	// ``K.''
	cased := false
	if symbol2 == 0 {
		upper := unicode.IsUpper(symbol1Rune)
		lower := unicode.IsLower(symbol1Rune)
		if upper || lower {
			symbol1Rune = unicode.ToLower(symbol1Rune)
			symbol2Rune = unicode.ToUpper(symbol1Rune)
			cased = true
		} else {
			symbol2     = symbol1
			symbol2Rune = symbol1Rune
		}
	}
	if symbol4 == 0 {
		upper := unicode.IsUpper(symbol3Rune)
		lower := unicode.IsLower(symbol3Rune)
		if upper || lower {
			symbol3Rune = unicode.ToLower(symbol3Rune)
			symbol4Rune = unicode.ToUpper(symbol3Rune)
			cased = true
		} else {
			symbol4     = symbol3
			symbol4Rune = symbol3Rune
		}
	}

	// PARAGRAPH 5
	//
	// The standard rules for obtaining a KeySym from a KeyPress event make
	// use of only the Group 1 and Group 2 KeySyms; no interpretation of/
	// other KeySyms in the list is given. Which group to use is determined
	// by the modifier state. Switching between groups is controlled by the
	// KeySym named MODE SWITCH, by attaching that KeySym to some KeyCode
	// and attaching that KeyCode to any one of the modifiers Mod1 through
	// Mod5. This modifier is called the group modifier. For any KeyCode,
	// Group 1 is used when the group modifier is off, and Group 2 is used
	// when the group modifier is on.
	modeSwitch := state & backend.modifierMasks.modeSwitch > 0
	if modeSwitch {
		symbol1     = symbol3
		symbol1Rune = symbol3Rune
		symbol2     = symbol4
		symbol2Rune = symbol4Rune
		
	}
	
	// PARAGRAPH 6
	// 
	// The Lock modifier is interpreted as CapsLock when the KeySym named
	// XK_Caps_Lock is attached to some KeyCode and that KeyCode is attached
	// to the Lock modifier. The Lock modifier is interpreted as ShiftLock
	// when the KeySym named XK_Shift_Lock is attached to some KeyCode and
	// that KeyCode is attached to the Lock modifier. If the Lock modifier
	// could be interpreted as both CapsLock and ShiftLock, the CapsLock
	// interpretation is used.
	shift :=
		state & xproto.ModMaskShift                > 0 ||
		state & backend.modifierMasks.shiftLock    > 0
	capsLock := state & backend.modifierMasks.capsLock > 0

	// PARAGRAPH 7
	//
	// The operation of keypad keys is controlled by the KeySym named
	// XK_Num_Lock, by attaching that KeySym to some KeyCode and attaching
	// that KeyCode to any one of the modifiers Mod1 through Mod5 . This
	// modifier is called the numlock modifier. The standard KeySyms with
	// the prefix ``XK_KP_'' in their name are called keypad KeySyms; these
	// are KeySyms with numeric value in the hexadecimal range 0xFF80 to
	// 0xFFBD inclusive. In addition, vendor-specific KeySyms in the
	// hexadecimal range 0x11000000 to 0x1100FFFF are also keypad KeySyms.
	numLock  := state & backend.modifierMasks.numLock  > 0
	
	// PARAGRAPH 8
	//
	// Within a group, the choice of KeySym is determined by applying the
	// first rule that is satisfied from the following list:
	var selectedKeysym xproto.Keysym
	var selectedRune   rune
	_, symbol2IsNumPad := keypadCodeTable[symbol2]
	switch {
	case numLock && symbol2IsNumPad:
		// The numlock modifier is on and the second KeySym is a keypad
		// KeySym. In this case, if the Shift modifier is on, or if the
		// Lock modifier is on and is interpreted as ShiftLock, then the
		// first KeySym is used, otherwise the second KeySym is used.
		if shift {
			selectedKeysym = symbol1
			selectedRune   = symbol1Rune
		} else {
			selectedKeysym = symbol2
			selectedRune   = symbol2Rune
		}

	case !shift && !capsLock:
		// The Shift and Lock modifiers are both off. In this case, the
		// first KeySym is used.
		selectedKeysym = symbol1
		selectedRune   = symbol1Rune

	case !shift && capsLock:
		// The Shift modifier is off, and the Lock modifier is on and is
		// interpreted as CapsLock. In this case, the first KeySym is
		// used, but if that KeySym is lowercase alphabetic, then the
		// corresponding uppercase KeySym is used instead.
		if cased && unicode.IsLower(symbol1Rune) {
			selectedRune = symbol2Rune
		} else {
			selectedKeysym = symbol1
			selectedRune   = symbol1Rune
		}

	case shift && capsLock:
		// The Shift modifier is on, and the Lock modifier is on and is
		// interpreted as CapsLock. In this case, the second KeySym is
		// used, but if that KeySym is lowercase alphabetic, then the
		// corresponding uppercase KeySym is used instead.
		if cased && unicode.IsLower(symbol2Rune) {
			selectedRune = unicode.ToUpper(symbol2Rune)
		} else {
			selectedKeysym = symbol2
			selectedRune   = symbol2Rune
		}

	case shift:
		// The Shift modifier is on, or the Lock modifier is on and is
		// interpreted as ShiftLock, or both. In this case, the second
		// KeySym is used.
		selectedKeysym = symbol2
		selectedRune   = symbol2Rune
	}

	////////////////////////////////////////////////////////////////
	// all of the below stuff is specific to tomo's button codes. //
	////////////////////////////////////////////////////////////////

	// look up in control code table
	var isControl bool
	button, isControl = buttonCodeTable[selectedKeysym]
	if isControl { return }

	// look up in keypad table
	button, numberPad = keypadCodeTable[selectedKeysym]
	if numberPad { return }

	// otherwise, use the rune
	button = tomo.Key(selectedRune)
	
	return
}

// keysymToRune takes in an X keysym and outputs a utf32 code point. This
// function does not and should not handle keypad keys, as those are handled
// by Backend.keycodeToButton.
func keysymToRune (keysym xproto.Keysym) (character rune) {
	// X keysyms like 0xFF.. or 0xFE.. are non-character keys. these cannot
	// be converted so we return a zero.
	if (keysym >> 8) == 0xFF || (keysym >> 8) == 0xFE {
		character = 0
		return
	}
	
	// some X keysyms have a single bit set to 1 here. i believe this is to
	// prevent conflicts with existing codes. if we mask it off we will get
	// a correct utf-32 code point.
	if keysym & 0xF000000 == 0x1000000 {
		character = rune(keysym & 0x0111111)
		return
	}

	// if none of these things happened, we can safely (i think) assume that
	// the keysym is an exact utf-32 code point.
	character = rune(keysym)
	return
}
