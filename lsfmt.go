/*
 *  Formatter that makes column display like a ls command
 */

package lsfmt

import (
	"fmt"
	"io"
	"math"
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

func (this Formatter) CalcColumns(items []string, vertical bool) (columns []int, err error) {
	var cands [][]int
	if this.width/this.space > len(items) {
		cands = make([][]int, len(items))
	} else {
		cands = make([][]int, this.width/this.space)
	}
	for i, _ := range cands {
		cands[i] = make([]int, i+1)
	}

	for i, _ := range items {
		for c, _ := range cands {
			var idx int
			if vertical {
				idx = i * (c + 1) / len(items)
			} else {
				idx = i % (c + 1)
			}

			w := stringWidth(items[i]) + this.space
			if cands[c][idx] < w {
				cands[c][idx] = w
			}

			if sum(cands[c])-this.space > this.width {
				cands = cands[:c]
				break
			}
		}
	}

	if len(cands) == 0 {
		longest := 0
		for _, s := range items {
			if longest < len(s) {
				longest = len(s)
			}
		}
		return nil, fmt.Errorf("terminal too narrow. this terminal has %d columns but longest string is %d characters.", this.width, longest)
	}

	columns = cands[len(cands)-1]
	columns[len(columns)-1] -= this.space
	return
}

func (this Formatter) PrintHorizontal(items []string) (columns []int, err error) {
	columns, err = this.CalcColumns(items, false)
	if err != nil {
		return nil, err
	}

	for i, _ := range items {
		fmt.Fprint(this.out, items[i])
		if i%len(columns) == len(columns)-1 {
			fmt.Fprintln(this.out)
		} else if i != len(items)-1 {
			space := columns[i%len(columns)] - stringWidth(items[i])
			fmt.Fprint(this.out, strings.Repeat(" ", space))
		}
	}
	if (len(items)-1)%len(columns) != len(columns)-1 {
		fmt.Fprintln(this.out)
	}

	return
}

func (this Formatter) PrintVertical(items []string) (columns []int, err error) {
	columns, err = this.CalcColumns(items, true)
	if err != nil {
		return nil, err
	}

	height := 1
	if len(columns) < len(items) {
		height = int(math.Ceil(float64(len(items)) / float64(len(columns))))
	}

	for r := 0; r < height; r++ {
		for c, _ := range columns {
			i := r + c*height
			if i >= len(items) {
				break
			} else {
				fmt.Fprint(this.out, items[i])
				if i+height < len(items) {
					space := columns[c] - stringWidth(items[i])
					fmt.Fprint(this.out, strings.Repeat(" ", space))
				}
			}
		}
		fmt.Fprintln(this.out)
	}

	return
}
