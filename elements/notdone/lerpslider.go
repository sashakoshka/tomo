package elements

// Numeric is a type constraint representing a number.
type Numeric interface {
	~float32 | ~float64 |
	~int     | ~int8    | ~int16  | ~int32  | ~int64  |
	~uint    | ~uint8   | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// LerpSlider is a slider that has a minimum and maximum value, and who's value
// can be any numeric type.
type LerpSlider[T Numeric] struct {
	*Slider
	min T
	max T
}

// NewLerpSlider creates a new LerpSlider with a minimum and maximum value. If
// vertical is set to true, the slider will be vertical instead of horizontal.
func NewLerpSlider[T Numeric] (min, max T, value T, vertical bool) (element *LerpSlider[T]) {
	if min > max {
		temp := max
		max = min
		min = temp
	}
	element = &LerpSlider[T] {
		Slider: NewSlider(0, vertical),
		min: min,
		max: max,
	}
	element.SetValue(value)
	return
}

// SetValue sets the slider's value.
func (element *LerpSlider[T]) SetValue (value T) {
	value -= element.min
	element.Slider.SetValue(float64(value) / float64(element.Range()))
}

// Value returns the slider's value.
func (element *LerpSlider[T]) Value () (value T) {
	return T (
		float64(element.Slider.Value()) * float64(element.Range())) +
		element.min
}

// Range returns the difference between the slider's maximum and minimum values.
func (element *LerpSlider[T]) Range () T {
	return element.max - element.min
}
