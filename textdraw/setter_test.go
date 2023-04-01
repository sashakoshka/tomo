package textdraw

import "image"
import "testing"
import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo/fixedutil"
import defaultfont "git.tebibyte.media/sashakoshka/tomo/default/font"

func TestSetterLength (test *testing.T) {
	// case 1
	text := []rune("The quick brown fox\njumped over the lazy dog.")
	setter := TypeSetter { }
	setter.SetText(text)
	setter.SetFace(defaultfont.FaceRegular)
	length := 0
	setter.For (func (i int, r rune, p fixed.Point26_6) bool {
		length ++
		return true
	})
	if length != len(text) {
		test.Fatalf (
			`setter rune count: %d, expected: %d`,
			length, len(text))
	}

	// case 2
	setter.SetMaxWidth(10)
	length = 0
	setter.For (func (i int, r rune, p fixed.Point26_6) bool {
		length ++
		return true
	})
	if length != len(text) {
		test.Fatalf (
			`setter rune count: %d, expected: %d`,
			length, len(text))
	}
}

func TestSetterBounds (test *testing.T) {
	setter := TypeSetter { }
	setter.SetText([]rune("The quick brown fox\njumped over the lazy dog."))
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

	testLargeRecHeight(test, 256)
	testLargeRecHeight(test, 100)
	testLargeRecHeight(test, 20)
	testLargeRecHeight(test, 400)
}

func testLargeRecHeight (test *testing.T, width int) {
	setter := TypeSetter { }
	setter.SetText([]rune(lipsum))
	setter.SetFace(defaultfont.FaceRegular)
	setter.SetMaxWidth(width)
	recHeight := setter.ReccomendedHeightFor(width)
	bounds := setter.LayoutBounds()

	if recHeight != bounds.Dy() {
		test.Fatalf (
			`setter bounds mismatch rec. height: %d, Dy: %d ` +
			`for width: %d`,
			recHeight, bounds.Dy(), width)
	}
}

func TestSetterPosition (test *testing.T) {
	setter := TypeSetter { }
	setter.SetText([]rune("The quick brown fox\njumped over the lazy dog."))
	setter.SetFace(defaultfont.FaceRegular)
	index := 20
	pos := fixedutil.RoundPt(setter.PositionAt(index))
	expect := image.Pt(0, 13)

	if pos != expect {
		test.Fatalf (
			`setter pos at %d: (%d, %d), expected: (%d, %d)`,
			index, pos.X, pos.Y, expect.X, expect.Y)
	}
}

func TestSetterIndex (test *testing.T) {
	setter := TypeSetter { }
	setter.SetText([]rune("The quick brown fox\njumped over the lazy dog."))
	setter.SetFace(defaultfont.FaceRegular)
	
	pos := fixed.P(3, 8)
	index := setter.AtPosition(pos)
	expect := 20
	if index != expect {
		test.Fatalf (
			`setter index at (%d, %d): %d, expected: %d`,
			pos.X.Round(), pos.Y.Round(), index, expect)
	}
	
	pos = fixed.P(-59, 230)
	index = setter.AtPosition(pos)
	expect = 20
	if index != expect {
		test.Fatalf (
			`setter index at (%d, %d): %d, expected: %d`,
			pos.X.Round(), pos.Y.Round(), index, expect)
	}
	
	pos = fixed.P(-500, -500)
	index = setter.AtPosition(pos)
	expect = 0
	if index != expect {
		test.Fatalf (
			`setter index at (%d, %d): %d, expected: %d`,
			pos.X.Round(), pos.Y.Round(), index, expect)
	}
	
	pos = fixed.P(500, -500)
	index = setter.AtPosition(pos)
	expect = 19
	if index != expect {
		test.Fatalf (
			`setter index at (%d, %d): %d, expected: %d`,
			pos.X.Round(), pos.Y.Round(), index, expect)
	}
	
	pos = fixed.P(500, 500)
	index = setter.AtPosition(pos)
	expect = setter.Length()
	if index != expect {
		test.Fatalf (
			`setter index at (%d, %d): %d, expected: %d`,
			pos.X.Round(), pos.Y.Round(), index, expect)
	}
}

const lipsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Fermentum et sollicitudin ac orci phasellus egestas tellus rutrum. Aliquam vestibulum morbi blandit cursus risus at ultrices mi. Gravida dictum fusce ut placerat. Cursus metus aliquam eleifend mi in nulla posuere. Sit amet nulla facilisi morbi tempus iaculis urna id. Amet volutpat consequat mauris nunc congue nisi vitae. Varius duis at consectetur lorem donec massa sapien faucibus et. Vitae elementum curabitur vitae nunc sed velit dignissim. In hac habitasse platea dictumst quisque sagittis purus. Enim nulla aliquet porttitor lacus luctus accumsan tortor. Lectus magna fringilla urna porttitor rhoncus dolor purus non.\n\nNon pulvinar neque laoreet suspendisse. Viverra adipiscing at in tellus integer. Vulputate dignissim suspendisse in est ante. Purus in mollis nunc sed id semper. In est ante in nibh mauris cursus. Risus pretium quam vulputate dignissim suspendisse in est. Blandit aliquam etiam erat velit scelerisque in dictum. Lectus quam id leo in. Odio tempor orci dapibus ultrices in iaculis. Pharetra sit amet aliquam id. Elit ut aliquam purus sit. Egestas dui id ornare arcu odio ut sem nulla pharetra. Massa tempor nec feugiat nisl pretium fusce id. Dui accumsan sit amet nulla facilisi morbi. A lacus vestibulum sed arcu non odio euismod. Nam libero justo laoreet sit amet cursus. Mattis rhoncus urna neque viverra justo nec. Mauris augue neque gravida in fermentum et sollicitudin ac. Vulputate mi sit amet mauris. Ut sem nulla pharetra diam sit amet."
