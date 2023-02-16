package textdraw

import "testing"
import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"

func TestSetterLength (test *testing.T) {
	text := []rune("The quick brown fox\njumped over the lazy dog.")
	setter := TypeSetter { }
	setter.SetText(text)
	setter.SetFace(defaultfont.FaceRegular)
	length := 0
	setter.For (func (i int, r rune, p fixed.Point26_6) bool {
		length ++
		return true
	})
	if length != len(text) - 1 {
		test.Fatalf (
			`setter rune count: %d, expected: %d`,
			length, len(text) - 1)
	}
}

func TestSetterBounds (test *testing.T) {
	text := []rune("The quick brown fox\njumped over the lazy dog.")
	setter := TypeSetter { }
	setter.SetText(text)
	setter.SetFace(defaultfont.FaceRegular)
	bounds := setter.LayoutBounds()
	
	expectDy := 13 * 2
	if expectDy != bounds.Dy() {
		test.Fatalf (
			`setter bounds Dy: %d, expected: %d`,
			bounds.Dy(), expectDy)
	}
	
	expectDx := 7 * 25
	if expectDx != bounds.Dx() {
		test.Fatalf (
			`setter bounds Dx: %d, expected: %d`,
			bounds.Dx(), expectDx)
	}
}
