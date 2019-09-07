package fs

import (
	"os"
	"path/filepath"
)

func DeleteFilesUnderDir(outputdir string) error {
	return filepath.Walk(outputdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if path == outputdir {
			return nil
		}
		if err := os.RemoveAll(path); err != nil {
			return nil
		} else {
			return err
		}
	})
}

