package gomega

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/types"

	"github.com/jlandowner/helm-chartsnap/pkg/snap"
)

var (
	shotCountMap = map[string]int{}
	trimSpace    = regexp.MustCompile(` +`)
)

// MatchSnapShot returns a Gomega matcher that compares the actual value with the snapshot file.
func MatchSnapShot() types.GomegaMatcher {

	testFile := ginkgo.CurrentSpecReport().FileName()
	path := filepath.Dir(testFile)
	file := filepath.Base(testFile)
	snapFile := filepath.Join(path, "__snapshots__", strings.TrimSuffix(file, ".go")+".snap")

	testLabel := ginkgo.CurrentSpecReport().FullText()
	testLabel = trimSpace.ReplaceAllString(testLabel, " ")

	count := shotCountMap[testLabel]
	count++
	shotCountMap[testLabel] = count
	snapId := fmt.Sprintf("%s %d", testLabel, count)

	return snap.SnapshotMatcher(snapFile, snap.WithSnapshotID(snapId))
}
