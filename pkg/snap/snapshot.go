package snap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
)

var (
	logger *slog.Logger
	mutex  sync.Mutex
)

func SetLogger(slogr *slog.Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	logger = slogr
}

func log() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	if logger == nil {
		logger = slog.Default()
	}
	return logger
}

func defaultDiffFunc(x, y string) string {
	var act interface{}
	if err := json.Unmarshal([]byte(x), &act); err != nil {
		return "Expected to match\n" + cmp.Diff(x, y)
	}
	var exp interface{}
	if err := json.Unmarshal([]byte(y), &exp); err != nil {
		return "Expected to match\n" + cmp.Diff(x, y)
	}
	return "Expected to JSON match\n" + cmp.Diff(exp, act)
}

// Option is a functional option for snapshotMatcher.
type Option func(m *snapshotMatcher)

func WithDiffFunc(f DiffFunc) Option {
	return func(m *snapshotMatcher) {
		m.diffFunc = f
	}
}

// WithSnapshotID is an option to specify the snapshot ID. If this option is set, the snapshot file is treated as a multi-snapshot file.
func WithSnapshotID(id string) Option {
	return func(m *snapshotMatcher) {
		m.snapID = id
	}
}

// SnapshotMatcher returns a matcher that compares the actual value with the snapshot file.
func SnapshotMatcher(snapFile string, options ...Option) *snapshotMatcher {
	m := &snapshotMatcher{
		snapFilePath: snapFile,
		diffFunc:     defaultDiffFunc,
	}

	for _, opt := range options {
		opt(m)
	}
	return m
}

type DiffFunc func(x, y string) string

type snapshotMatcher struct {
	snapFilePath   string
	snapID         string
	expectedString string
	actualString   string
	diffFunc       DiffFunc
}

func (m *snapshotMatcher) Match(actual interface{}) (success bool, err error) {
	// prepare actual
	switch v := actual.(type) {
	case string:
		m.actualString = v
	case []byte:
		m.actualString = string(v)
	default:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(true)
		enc.SetIndent("", "  ")
		if err := enc.Encode(actual); err != nil {
			return false, fmt.Errorf("json encode error: %w", err)
		}
		m.actualString = buf.String()
	}

	// prepare expected from snapshot
	snap, err := m.readSnapshot()
	if errors.Is(err, afero.ErrFileNotFound) {
		// take new snapshot
		err = m.writeSnapshot([]byte(m.actualString))
		if err == nil {
			return true, nil
		}
	}
	if err != nil {
		return false, fmt.Errorf("preparing expected from snapshot error: %w", err)
	}
	m.expectedString = string(snap)

	return m.actualString == m.expectedString, nil
}

// FailureMessage returns a string that describes the failure of the matcher.
// actual must be always nil because the actual value is already parsed and stored in the matcher.
func (m *snapshotMatcher) FailureMessage(actual interface{}) (message string) {
	return m.diffFunc(m.expectedString, m.actualString) + "\n"
}

// NegatedFailureMessage returns a string that describes the failure of the negated matcher.
// actual must be always nil because the actual value is already parsed and stored in the matcher.
func (m *snapshotMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.diffFunc(m.expectedString, m.actualString) + "\n"
}

func (m *snapshotMatcher) readSnapshot() ([]byte, error) {
	raw, err := ReadFile(m.snapFilePath)
	if err != nil {
		return nil, err
	}
	if m.snapID == "" {
		return raw, nil
	}

	log().Debug("read multi snapshot", "snapFilePath", m.snapFilePath, "snapID", m.snapID)
	snaps, err := DecodeMultiSnapshots(raw)
	if err != nil {
		return nil, err
	}
	if snap, ok := snaps[m.snapID]; ok {
		s := snap.Snapshot.(string)
		return []byte(s), nil
	} else {
		return nil, afero.ErrFileNotFound
	}
}

func (m *snapshotMatcher) writeSnapshot(snapFileData []byte) error {
	if m.snapID == "" {
		return WriteFile(m.snapFilePath, snapFileData)

	}

	log().Debug("write multi snapshot", "snapFilePath", m.snapFilePath, "snapID", m.snapID)
	var snaps multiSnap
	raw, err := ReadFile(m.snapFilePath)
	if err == nil {
		snaps, err = DecodeMultiSnapshots(raw)
		if err != nil {
			return fmt.Errorf("decode snapshot file error: %w", err)
		}
	}
	if err != nil {
		if errors.Is(err, afero.ErrFileNotFound) {
			snaps = make(multiSnap)
		} else {
			return fmt.Errorf("read snapshot file error: %w", err)
		}
	}

	snaps[m.snapID] = snap{Snapshot: string(snapFileData)}

	data, err := EncodeMultiSnapshots(snaps)
	if err != nil {
		return err
	}
	return WriteFile(m.snapFilePath, data)
}

type multiSnap map[string]snap

type snap struct {
	Snapshot interface{} `toml:"SnapShot,multiline,omitempty"`
}

func IsMultiSnapshots(filePath string) bool {
	raw, err := ReadFile(filePath)
	if err != nil {
		return false
	}
	_, err = DecodeMultiSnapshots(raw)
	return err == nil
}

func DecodeMultiSnapshots(raw []byte) (multiSnap, error) {
	snaps := make(multiSnap)
	err := toml.NewDecoder(bytes.NewReader(raw)).Decode(&snaps)
	if err != nil {
		return nil, fmt.Errorf("toml decode error: %w", err)
	}
	if len(snaps) == 0 {
		return snaps, nil
	}
	return snaps, nil
}

func EncodeMultiSnapshots(snaps multiSnap) ([]byte, error) {
	buf := bytes.Buffer{}
	tomlEnc := toml.NewEncoder(&buf)
	tomlEnc.SetArraysMultiline(true)
	if err := tomlEnc.Encode(snaps); err != nil {
		return nil, fmt.Errorf("toml encode error: %w", err)
	}
	return buf.Bytes(), nil
}
