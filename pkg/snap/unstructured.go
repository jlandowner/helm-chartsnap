package snap

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aryann/difflib"
	"github.com/fatih/color"
	gomegatypes "github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	unstructutils "github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

func UnstructuredSnapShotMatcher(snapFile string, snapId string, diffOpts ...DiffOptions) *snapShotMatcher {
	o := mergeDiffOpts(diffOpts)

	return &snapShotMatcher{
		snapFilePath: snapFile,
		snapId:       snapId,
		fs:           cacheFs,
		diffFunc:     UnstructuredSnapshotDiff,
		diffOptions:  o,
	}
}

func UnstructuredMatch(matcher gomegatypes.GomegaMatcher, manifests []unstructured.Unstructured) (success bool, err error) {
	res, err := unstructutils.Encode(manifests)
	if err != nil {
		return false, fmt.Errorf("failed to encode manifests: %w", err)
	}
	return matcher.Match(string(res))
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

func UnstructuredSnapshotDiff(x, y string, o DiffOptions) string {
	divExp := regexp.MustCompile(`^  - object:$`)
	diffs := difflib.Diff(strings.Split(x, "\n"), strings.Split(y, "\n"))

	var (
		sb             strings.Builder
		isDiffSequence bool
		currentKind    string
		currentName    string
	)

	for i, v := range diffs {
		if o.ContextLineN() < 1 {
			// all records
			sb.WriteString(diffString(v))
			continue
		}

		if divExp.Match([]byte(v.String())) {
			isDiffSequence = false
			currentKind, currentName = findKind(diffs[i:]), findName(diffs[i:])
			Log().Debug("div match", "kind", currentKind, "name", currentName, "index", i)
		}

		if v.Delta != difflib.Common {
			isDiffSequence = true

			// if first diff, add a header and previous lines
			if i > 0 && diffs[i-1].Delta == difflib.Common {
				// header
				sb.WriteString(color.New(color.FgCyan).Sprintf("--- kind=%s name=%s line=%d\n", currentKind, currentName, i))

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
