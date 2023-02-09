package music

import "math"

var semitone = math.Pow(2, 1.0 / 12.0)

// Tuning is an interface representing a tuning.
type Tuning interface {
	// Tune returns the frequency of a given note in Hz.
	Tune (Note) float64
}

// EqualTemparment implements twelve-tone equal temparment.
type EqualTemparment struct { A4 float64 }

// Tune returns the EqualTemparment frequency of a given note in Hz.
func (tuning EqualTemparment) Tune (note Note) float64 {
	return tuning.A4 * math.Pow(semitone, float64(note - NoteA4))
}

// Octave represents a MIDI octave.
type Octave int

// Note returns the note at the specified scale degree in the chromatic scale.
func (octave Octave) Note (degree int) Note {
	return Note(int(octave + 1) * 12 + degree)
}

// Note represents a MIDI note.
type Note int

const (
	NoteC0 Note = iota
	NoteDb0
	NoteD0
	NoteEb0
	NoteE0
	NoteF0
	NoteGb0
	NoteG0
	NoteAb0
	NoteA0
	NoteBb0
	NoteB0

	NoteA4 Note = 69
)

// Octave returns the octave of the note
func (note Note) Octave () int {
	return int(note / 12 - 1)
}

// Degree returns the scale degree of the note in the chromatic scale.
func (note Note) Degree () int {
	mod := note % 12
	if mod < 0 { mod += 12 }
	return int(mod)
}

// IsSharp returns whether or not the note is a sharp.
func (note Note) IsSharp () bool {
	degree := note.Degree()
	return degree == 1 ||
		degree == 3 ||
		degree == 6 ||
		degree == 8 ||
		degree == 10
}
