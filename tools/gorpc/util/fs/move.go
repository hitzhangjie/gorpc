package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveFileToFolder move file or folder `src` to folder `dest`.
func Move(src, dest string) error {

	var (
		inf os.FileInfo
		err error
	)

	if inf, err = os.Lstat(dest); err != nil {
		return err
	}

	// move src to existed folder dest
	if !inf.IsDir() {
		return fmt.Errorf("dest:%s not folder", dest)
	}

	_, fname := filepath.Split(src)
	target := filepath.Join(dest, fname)

	// keep behavior consistent with bash `mv` command
	if _, err = os.Lstat(target); os.IsNotExist(err) {
		return os.Rename(src, filepath.Join(dest, target))
	}

	if err = os.RemoveAll(target); err != nil {
		return err
	}
	return os.Rename(src, filepath.Join(dest, target))
}
