package textdraw

import "testing"
import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"

func TestDoWord (test *testing.T) {
	text := []rune("The quick brown fox")
	word, remaining := DoWord(text, defaultfont.FaceRegular)
	
	expect := "quick brown fox"
	if string(remaining) != expect {
		test.Fatalf (
			`text: "%s", remaining: "%s" expected: "%s"`,
			string(text), string(remaining), expect)
	}

	if len(word.Runes) != 4 {
		test.Fatalf(`wrong rune length %d`, len(word.Runes))
	}

	if word.FirstRune() != 'T' {
		test.Fatalf(`wrong first rune %s`, string(word.FirstRune()))
	}

	if word.LastRune() != ' ' {
		test.Fatalf(`wrong last rune %s`, string(word.FirstRune()))
	}
}

func TestDoLine (test *testing.T) {
	// case 1
	text := []rune("The quick brown fox\njumped over the lazy dog")
	line, remaining := DoLine(text, defaultfont.FaceRegular, 0)
	
	expect := "jumped over the lazy dog"
	if string(remaining) != expect {
		test.Fatalf (
			`text: "%s", remaining: "%s" expected: "%s"`,
			string(text), string(remaining), expect)
	}

	if len(line.Words) != 4 {
		test.Fatalf(`wrong word count %d`, len(line.Words))
	}
	
	if !line.BreakAfter {
		test.Fatalf(`did not set BreakAfter to true`)
	}

	// case 2
	text = []rune("jumped over the lazy dog")
	line, remaining = DoLine(text, defaultfont.FaceRegular, 0)
	
	expect = ""
	if string(remaining) != expect {
		test.Fatalf (
			`text: "%s", remaining: "%s" expected: "%s"`,
			string(text), string(remaining), expect)
	}

	if len(line.Words) != 5 {
		test.Fatalf(`wrong word count %d`, len(line.Words))
	}
	
	if line.BreakAfter {
		test.Fatalf(`did not set BreakAfter to false`)
	}
	
	// case 3
	text = []rune("jumped over the lazy dog")
	line, remaining = DoLine(text, defaultfont.FaceRegular, fixed.I(7 * 12))
	
	expect = "the lazy dog"
	if string(remaining) != expect {
		test.Fatalf (
			`text: "%s", remaining: "%s" expected: "%s"`,
			string(text), string(remaining), expect)
	}

	if len(line.Words) != 2 {
		test.Fatalf(`wrong word count %d`, len(line.Words))
	}
	
	if line.BreakAfter {
		test.Fatalf(`did not set BreakAfter to false`)
	}
}
