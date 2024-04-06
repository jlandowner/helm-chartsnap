package snap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

var (
	cacheFs = afero.NewCacheOnReadFs(
		afero.NewOsFs(),
		afero.NewMemMapFs(),
		time.Minute,
	)
)

func WriteFile(path string, data []byte) error {
	log().Debug("write file by cacheFs", "path", path)

	if err := cacheFs.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("create shopshot directory error: %w", err)
	}
	file, err := cacheFs.Create(path)
	if err != nil {
		return fmt.Errorf("open shopshot file error: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write shopshot file error: %w", err)
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	log().Debug("read file by cacheFs", "path", path)

	exists, err := afero.Exists(cacheFs, path)
	if err != nil {
		return nil, fmt.Errorf("file check error: %w", err)
	}
	if !exists {
		return nil, afero.ErrFileNotFound
	}

	file, err := cacheFs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}
	return b, nil
}

func RemoveFile(path string) error {
	log().Debug("remove file by cacheFs", "path", path)
	return cacheFs.Remove(path)
}
