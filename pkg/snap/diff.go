package snap

import (
	"fmt"
	"strings"

	"github.com/aryann/difflib"
	"github.com/fatih/color"
)

type DiffOptions struct {
	DiffContextLineN int
}

func (o *DiffOptions) ContextLineN() int {
	if o.DiffContextLineN < 0 {
		return 0
	}
	return o.DiffContextLineN
}

func WithDiffContextLineN(n int) DiffOptions {
	return DiffOptions{DiffContextLineN: n}
}

func mergeDiffOpts(opts []DiffOptions) DiffOptions {
	var merged DiffOptions
	for _, v := range opts {
		if v.DiffContextLineN > merged.DiffContextLineN {
			merged.DiffContextLineN = v.DiffContextLineN
		}
	}
	return merged
}

func Diff(x, y string, o DiffOptions) string {
	diffs := difflib.Diff(strings.Split(x, "\n"), strings.Split(y, "\n"))

	var (
		sb             strings.Builder
		isDiffSequence bool
	)

	for i, v := range diffs {
		if o.ContextLineN() < 1 {
			// all records
			sb.WriteString(diffString(v))
			continue
		}

		if v.Delta != difflib.Common {
			isDiffSequence = true

			// if first diff, add a header and previous lines
			if i > 0 && diffs[i-1].Delta == difflib.Common {
				// header
				sb.WriteString(color.New(color.FgCyan).Sprintf("--- line=%d\n", i))

				// previous lines
				for j := intInRange(0, len(diffs), i-o.DiffContextLineN); j < i; j++ {
					sb.WriteString(fmt.Sprintf("%s\n", diffs[j]))
				}
			}
			sb.WriteString(diffString(v))
		} else {
			if isDiffSequence {
				isDiffSequence = false

				// subsequent lines
				for j := i; j < intInRange(0, len(diffs), i+o.DiffContextLineN); j++ {
					sb.WriteString(fmt.Sprintf("%s\n", diffs[j]))
				}
				// divider
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}

func intInRange(min, max, v int) int {
	if v >= min && v <= max {
		return v
	} else if v < min {
		return min
	} else {
		return max
	}
}

func diffString(d difflib.DiffRecord) string {
	switch d.Delta {
	case difflib.LeftOnly:
		return color.New(color.FgRed).Sprintf("%s\n", d)
	case difflib.RightOnly:
		return color.New(color.FgGreen).Sprintf("%s\n", d)
	default:
		return fmt.Sprintf("%s\n", d)
	}
}
