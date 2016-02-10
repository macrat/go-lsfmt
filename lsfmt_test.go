package lsfmt

import (
	"bytes"
	"testing"
)

func TestHorizontalASCII(t *testing.T) {
	items := []string{
		"hello",
		"world",
		"this",
		"is",
		"a",
		"test",
	}
	oneColumn := `hello
world
this
is
a
test
`
	twoColumn := `hello  world
this   is
a      test
`
	fourColumn := `hello world this is
a     test
`

	buf := bytes.NewBuffer(*new([]byte))
	fmt := NewFormatterWriter(buf, 4)

	cols, err := fmt.PrintHorizontal(items)
	if err == nil {
		t.Errorf("gave too narrow terminal but succeed printing with %d columns.", len(cols))
	}

	buf.Reset()
	fmt.width = 11
	cols, err = fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to printing in 11 columns terminal: %s", err)
	}
	if len(cols) != 1 {
		t.Errorf("expected columns is 1 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != oneColumn {
		t.Errorf("unexpected output in 11 columns terminal.")
	}

	buf.Reset()
	fmt.width = 12
	cols, err = fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to printing in 12 columns terminal: %s", err)
	}
	if len(cols) != 2 {
		t.Errorf("expected columns is 2 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != twoColumn {
		t.Errorf("unexpected output in 12 columns terminal.")
	}

	buf.Reset()
	fmt.width = 19
	fmt.space = 1
	cols, err = fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to printing in 19 columns terminal: %s", err)
	}
	if len(cols) != 4 {
		t.Errorf("expected columns is 4 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != fourColumn {
		t.Errorf("unexpected output in 19 columns terminal.")
	}
}

func TestHorizontalJapanese(t *testing.T) {
	items := []string{
		"あいう",
		"「」",
		"abc",
		"アイウ",
		"ｱｲｳ",
		"def",
		"漢字",
		"、。",
		"hij",
		"klm",
		"nop",
		"qrs",
	}

	expected := `あいう  「」  abc
アイウ  ｱｲｳ   def
漢字    、。  hij
klm     nop   qrs
`

	buf := bytes.NewBuffer(*new([]byte))
	fmt := NewFormatterWriter(buf, 20)

	cols, err := fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to printing: %s", err)
	}
	if len(cols) != 3 {
		t.Errorf("expected columns is 3 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("unexpected output.")
	}
}
func TestVertical(t *testing.T) {
	items := []string{
		"hello",
		"world",
		"this",
		"is",
		"test",
	}
	oneColumn := `hello
world
this
is
test
`
	twoColumn := `hello  is
world  test
this
`
	fourColumn := `hello this test
world is
`

	buf := bytes.NewBuffer(*new([]byte))
	fmt := NewFormatterWriter(buf, 7)

	cols, err := fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to printing in 8 columns terminal: %s", err)
	}
	if len(cols) != 1 {
		t.Errorf("expected columns is 1 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != oneColumn {
		t.Errorf("unexpected output in 8 columns terminal.")
	}

	buf.Reset()
	fmt.width = 11
	cols, err = fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to printing in 11 columns terminal: %s", err)
	}
	if len(cols) != 2 {
		t.Errorf("expected columns is 2 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != twoColumn {
		t.Errorf("unexpected output in 11 columns terminal.")
	}

	buf.Reset()
	fmt.width = 20
	fmt.space = 1
	cols, err = fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to printing in 20 columns terminal: %s", err)
	}
	if len(cols) != 4 {
		t.Errorf("expected columns is 4 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != fourColumn {
		t.Errorf("unexpected output in 20 columns terminal.")
	}
}
