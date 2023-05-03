package elements

import "tomo"

// Numeric is a type constraint representing a number.
type Numeric interface {
	~float32 | ~float64 |
	~int     | ~int8    | ~int16  | ~int32  | ~int64  |
	~uint    | ~uint8   | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// LerpSlider is a slider that has a minimum and maximum value, and who's value
// can be any numeric type.
type LerpSlider[T Numeric] struct {
	slider
	min T
	max T
}

// NewVLerpSlider creates a new horizontal LerpSlider with a minimum and maximum
// value.
func NewVLerpSlider[T Numeric] (min, max T, value T) (element *LerpSlider[T]) {
	element = NewHLerpSlider(min, max, value)
	element.vertical = true
	return
}

// NewHLerpSlider creates a new horizontal LerpSlider with a minimum and maximum
// value.
func NewHLerpSlider[T Numeric] (min, max T, value T) (element *LerpSlider[T]) {
	if min > max { min, max = max, min }
	element = &LerpSlider[T] {
		min: min,
		max: max,
	}
	element.entity = tomo.GetBackend().NewEntity(element)
	element.construct()
	element.SetValue(value)
	return
}

// SetValue sets the slider's value.
func (element *LerpSlider[T]) SetValue (value T) {
	value -= element.min
	element.slider.SetValue(float64(value) / float64(element.Range()))
}

// Value returns the slider's value.
func (element *LerpSlider[T]) Value () (value T) {
	return T (
		float64(element.slider.Value()) * float64(element.Range())) +
		element.min
}

// Range returns the difference between the slider's maximum and minimum values.
func (element *LerpSlider[T]) Range () T {
	return element.max - element.min
}
