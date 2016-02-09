/*
 *  Formatter that makes column display like a ls command
 */

package lsfmt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/text/width"
)

func stringWidth(str string) (w int) {
	for _, c := range str {
		kind := width.LookupRune(c).Kind()
		if kind == width.EastAsianWide || kind == width.EastAsianFullwidth {
			w += 2
		} else {
			w += 1
		}
	}
	return
}

func sum(xs []int) (s int) {
	for _, x := range xs {
		s += x
	}
	return
}

type Formatter struct {
	out   io.Writer
	width int // width of one line.
	space int // number of spaces between items.
}

func NewFormatterWriter(writer io.Writer, width int) (formatter Formatter) {
	return Formatter{out: writer, width: width, space: 2}
}

func NewFormatterFile(file *os.File) (formatter Formatter, err error) {
	width, _, err := terminal.GetSize(int(file.Fd()))
	formatter = Formatter{out: file, width: width, space: 2}
	return
}

func (this Formatter) Print(items []string) (columns int, err error) {
	var cands [][]int
	for i := 0; i < this.width/this.space; i++ {
		cands = append(cands, *new([]int))
	}

	for i := 0; i < len(items); i++ {
		for c := 0; c < len(cands); c++ {
			idx := i % (c + 1)

			w := stringWidth(items[i])
			if idx < c {
				w += this.space
			}

			if len(cands[c]) <= idx {
				cands[c] = append(cands[c], w)
			} else if cands[c][idx] < w {
				cands[c][idx] = w
			}

			if sum(cands[c][:]) > this.width {
				cands = cands[:c]
				break
			}
		}
	}

	columns = len(cands)

	if columns == 0 {
		longest := 0
		for _, s := range items {
			if longest < len(s) {
				longest = len(s)
			}
		}
		err = fmt.Errorf("terminal too narrow. this terminal has %d columns but longest string is %d characters.", this.width, longest)
		return
	}

	for i := 0; i < len(items); i++ {
		fmt.Fprint(this.out, items[i])
		if i%columns == columns-1 {
			fmt.Fprintln(this.out)
		} else if i != len(items)-1 {
			fmt.Fprint(this.out, strings.Repeat(" ", cands[len(cands)-1][i%columns]-stringWidth(items[i])))
		}
	}
	if (len(items)-1)%columns != columns-1 {
		fmt.Fprintln(this.out)
	}

	return
}
