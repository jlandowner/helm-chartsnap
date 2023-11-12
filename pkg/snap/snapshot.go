package snap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
)

var log *slog.Logger

func SetLogger(slogr *slog.Logger) {
	log = slogr
}

func Log() *slog.Logger {
	if log == nil {
		log = slog.Default()
	}
	return log
}

var (
	shotCountMap = map[string]int{}
	cacheFs      = afero.NewCacheOnReadFs(
		afero.NewOsFs(),
		afero.NewMemMapFs(),
		time.Minute,
	)
	trimSpace = regexp.MustCompile(` +`)
)

func MatchSnapShot(options ...Option) types.GomegaMatcher {

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

	return SnapShotMatcher(snapFile, snapId)
}

func SnapShotMatcher(snapFile string, snapId string, diffOpts ...DiffOptions) *snapShotMatcher {
	o := mergeDiffOpts(diffOpts)

	return &snapShotMatcher{
		snapFilePath: snapFile,
		snapId:       snapId,
		fs:           cacheFs,
		diffFunc:     Diff,
		diffOptions:  o,
	}
}

type Option func(m *snapShotMatcher)

type snapShotMatcher struct {
	snapFilePath string
	snapId       string
	fs           afero.Fs
	expectedJson string
	actualJson   string
	diffFunc     func(x, y string, opts DiffOptions) string
	diffOptions  DiffOptions
}

func (m *snapShotMatcher) Match(actual interface{}) (success bool, err error) {

	switch v := actual.(type) {
	case string:
		m.actualJson = v
	default:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(true)
		enc.SetIndent("", "  ")
		if err := enc.Encode(actual); err != nil {
			return false, fmt.Errorf("json encode error: %w", err)
		}
		m.actualJson = buf.String()
	}

	snap, err := m.ReadSnapShot()
	if errors.Is(err, afero.ErrFileNotFound) {
		err = m.WriteSnapShot([]byte(m.actualJson))
		if err == nil {
			return true, nil
		}
	}
	if err != nil {
		return false, err
	}
	m.expectedJson = *snap

	return m.actualJson == m.expectedJson, nil
}

func (m *snapShotMatcher) FailureMessage(actual interface{}) (message string) {
	return "Expected to match\n" + m.diffFunc(m.expectedJson, m.actualJson, m.diffOptions) + "\n"
}

func (m *snapShotMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	var act interface{}
	if err := json.Unmarshal([]byte(m.actualJson), &act); err != nil {
		return fmt.Errorf("json decode error: %w", err).Error()
	}
	var exp interface{}
	if err := json.Unmarshal([]byte(m.expectedJson), &exp); err != nil {
		return fmt.Errorf("json decode error: %w", err).Error()
	}
	return "Expected not to match\n" + cmp.Diff(exp, act) + "\n"
}

//-----------------------------------------------------------

type data struct {
	SnapShot interface{} `toml:"SnapShot,multiline,omitempty"`
}

func (m *snapShotMatcher) ReadSnapShot() (*string, error) {

	snapFileData, err := m.readSnapFileData()
	if err != nil {
		return nil, err
	}

	if snap, ok := (*snapFileData)[m.snapId]; ok {
		a := snap.SnapShot.(string)
		return &a, nil
	} else {
		return nil, afero.ErrFileNotFound
	}
}

func (m *snapShotMatcher) WriteSnapShot(snap []byte) error {

	snapFileData, err := m.readSnapFileData()
	if err != nil {
		return err
	}
	(*snapFileData)[m.snapId] = data{SnapShot: string(snap)}

	if err := m.writeSnapFileData(snapFileData); err != nil {
		return err
	}
	return nil
}

func (m *snapShotMatcher) readSnapFileData() (*map[string]data, error) {
	exists, err := afero.Exists(m.fs, m.snapFilePath)
	if err != nil {
		return nil, fmt.Errorf("file check error: %w", err)
	}
	if !exists {
		return &map[string]data{}, nil
	}

	file, err := m.fs.Open(m.snapFilePath)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}

	defer file.Close()

	var datas map[string]data
	err = toml.NewDecoder(file).Decode(&datas)
	if err != nil {
		return nil, fmt.Errorf("toml decode error: %w", err)
	}
	if len(datas) == 0 {
		return &map[string]data{}, nil
	}
	return &datas, nil
}

func (m *snapShotMatcher) writeSnapFileData(snapFileData *map[string]data) error {
	if err := m.fs.MkdirAll(filepath.Dir(m.snapFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("create snapfile directory error: %w", err)
	}
	file, err := m.fs.Create(m.snapFilePath)
	if err != nil {
		return fmt.Errorf("open snapfile error: %w", err)
	}
	defer file.Close()

	enc := toml.NewEncoder(file)
	enc.SetArraysMultiline(true)

	if err := enc.Encode(snapFileData); err != nil {
		return fmt.Errorf("toml encode error: %w", err)
	}
	return nil
}
