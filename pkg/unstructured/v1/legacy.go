package unstructured

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/aryann/difflib"
	"github.com/fatih/color"
	yaml "go.yaml.in/yaml/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Encode encoding legacy formatted yaml
func Encode(arr []metaV1.Unstructured) ([]byte, error) {
	sort.SliceStable(arr, func(i, j int) bool {
		if arr[i].GetAPIVersion() != arr[j].GetAPIVersion() {
			return arr[i].GetAPIVersion() < arr[j].GetAPIVersion()
		}
		if arr[i].GetKind() != arr[j].GetKind() {
			return arr[i].GetKind() < arr[j].GetKind()
		}
		return arr[i].GetName() < arr[j].GetName()
	})
	return yaml.Marshal(arr)
}

// extract kind value
func findKind(diffs []difflib.DiffRecord) string {
	kindExp := regexp.MustCompile(`^      kind: (.+)$`)
	for i := 0; i < len(diffs); i++ {
		kindMatch := kindExp.FindStringSubmatch(diffs[i].String())
		if len(kindMatch) > 0 {
			return kindMatch[1]
		}
	}
	return ""
}

// extract name value
func findName(diffs []difflib.DiffRecord) string {
	metaExp := regexp.MustCompile(`^      metadata:$`)
	nameExp := regexp.MustCompile(`^          name: (.+)$`)
	for i := 0; i < len(diffs); i++ {
		if metaExp.Match([]byte(diffs[i].String())) {
			for j := i + 1; j < len(diffs)-i; j++ {
				nameMatch := nameExp.FindStringSubmatch(diffs[j].String())
				if len(nameMatch) > 0 {
					return nameMatch[1]
				}
			}
			return ""
		}
	}
	return ""
}

type DiffOptions struct {
	ContextLineN int
}

func (o *DiffOptions) Diff(x, y string) string {
	divExp := regexp.MustCompile(`^  - object:$`)
	diffs := difflib.Diff(strings.Split(x, "\n"), strings.Split(y, "\n"))

	var (
		sb             strings.Builder
		isDiffSequence bool
		currentKind    string
		currentName    string
	)

	for i, v := range diffs {
		if o.ContextLineN < 1 {
			// all records
			sb.WriteString(diffString(v))
			continue
		}

		if divExp.Match([]byte(v.String())) {
			isDiffSequence = false
			currentKind, currentName = findKind(diffs[i:]), findName(diffs[i:])
		}

		if v.Delta != difflib.Common {
			isDiffSequence = true

			// if first diff, add a header and previous lines
			if i > 0 && diffs[i-1].Delta == difflib.Common {
				// header
				sb.WriteString(color.New(color.FgCyan).Sprintf("@@ KIND=%s NAME=%s LINE=%d\n", currentKind, currentName, i))

				// previous lines
				for j := intInRange(0, len(diffs), i-o.ContextLineN); j < i; j++ {
					sb.WriteString(fmt.Sprintf("%s\n", diffs[j]))
				}
			}
			sb.WriteString(diffString(v))

		} else {
			if isDiffSequence {
				isDiffSequence = false

				// subsequent lines
				for j := i; j < intInRange(0, len(diffs), i+o.ContextLineN); j++ {
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
