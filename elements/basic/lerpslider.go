package basicElements

// Numeric is a type constraint representing a number.
type Numeric interface {
	~float32 | ~float64 |
	~int     | ~int8    | ~int16  | ~int32  | ~int64  |
	~uint    | ~uint8   | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type LerpSlider[T Numeric] struct {
	*Slider
	min T
	max T
}

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

func (element *LerpSlider[T]) SetValue (value T) {
	value -= element.min
	element.Slider.SetValue(float64(value) / float64(element.Range()))
}

func (element *LerpSlider[T]) Value () (value T) {
	return T (
		float64(element.Slider.Value()) * float64(element.Range())) +
		element.min
}

func (element *LerpSlider[T]) Range () T {
	return element.max - element.min
}
