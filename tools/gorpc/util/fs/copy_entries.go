package fs

import (
	"os"
	"path/filepath"
)

func CopyFileUnderDir(srcDir string, destDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, inerr error) error {
		if inerr != nil {
			return inerr
		}
		if info.IsDir() {
			return nil
		}
		if err := Copy(path, filepath.Join(destDir, info.Name())); err != nil {
			return err
		}
		return nil
	})
}
