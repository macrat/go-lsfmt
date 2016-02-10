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
		t.Errorf("failed to printing in %d columns terminal: %s", fmt.width, err)
	}
	if len(cols) != 1 {
		t.Errorf("expected columns is 1 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != oneColumn {
		t.Errorf("unexpected output in %d columns terminal.", fmt.width)
	}

	buf.Reset()
	fmt.width = 12
	cols, err = fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to printing in %d columns terminal: %s", fmt.width, err)
	}
	if len(cols) != 2 {
		t.Errorf("expected columns is 2 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != twoColumn {
		t.Errorf("unexpected output in %d columns terminal.", fmt.width)
	}

	buf.Reset()
	fmt.width = 19
	fmt.space = 1
	cols, err = fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to printing in %d columns terminal: %s", fmt.width, err)
	}
	if len(cols) != 4 {
		t.Errorf("expected columns is 4 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != fourColumn {
		t.Errorf("unexpected output in %d columns terminal.", fmt.width)
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
		"a",
		"test",
		"of",
		"vertical",
		"printing",
	}
	oneColumn := `hello
world
this
is
a
test
of
vertical
printing
`
	twoColumn := `hello test
world of
this  vertical
is    printing
a
`
	threeColumn := `hello  is    of
world  a     vertical
this   test  printing
`

	buf := bytes.NewBuffer(*new([]byte))
	fmt := NewFormatterWriter(buf, 10)

	cols, err := fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to printing in %d columns terminal: %s", fmt.width, err)
	}
	if len(cols) != 1 {
		t.Errorf("expected columns is 1 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != oneColumn {
		t.Errorf("unexpected output in %d columns terminal.", fmt.width)
	}

	buf.Reset()
	fmt.width = 15
	fmt.space = 1
	cols, err = fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to printing in %d columns terminal: %s", fmt.width, err)
	}
	if len(cols) != 2 {
		t.Errorf("expected columns is 2 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != twoColumn {
		t.Errorf("unexpected output in %d columns terminal.", fmt.width)
	}

	buf.Reset()
	fmt.width = 24
	fmt.space = 2
	cols, err = fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to printing in %d columns terminal: %s", fmt.width, err)
	}
	if len(cols) != 3 {
		t.Errorf("expected columns is 3 but got %d columns output.", len(cols))
	}
	if string(buf.Bytes()) != threeColumn {
		t.Errorf("unexpected output in %d columns terminal.", fmt.width)
	}
}

func TestShort(t *testing.T) {
	items := []string{
		"hello",
		"world",
	}
	expected := "hello  world\n"

	buf := bytes.NewBuffer(*new([]byte))
	fmt := NewFormatterWriter(buf, 80)

	cols, err := fmt.PrintHorizontal(items)
	if err != nil {
		t.Errorf("failed to horizontal printing: %s", err)
	}
	if len(cols) != 2 {
		t.Errorf("expected columns is 2 but got %d columns horizontal output.", len(cols))
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("unexpected output in horizontal printing.")
	}

	buf.Reset()
	cols, err = fmt.PrintVertical(items)
	if err != nil {
		t.Errorf("failed to vertical printing: %s", err)
	}
	if len(cols) != 2 {
		t.Errorf("expected columns is 2 but got %d columns output in vertical printing.", len(cols))
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("unexpected output in vertical printing.")
	}
}
