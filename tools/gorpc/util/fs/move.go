package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveFileToFolder move file to folder with name filename `name`.
func Move(src, dest string) error {

	inf, err := os.Lstat(dest)

	// move src to existed folder/file dest
	if err == nil {
		if inf.IsDir() {
			_, file := filepath.Split(src)
			return os.Rename(src, filepath.Join(dest, file))
		} else {
			return fmt.Errorf("dest not folder, %s", dest)
		}
	}

	// move src to unexisted folder/file dest
	dir, _ := filepath.Split(dest)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Rename(src, dest)
	if err != nil {
		return err
	}

	return nil
}

