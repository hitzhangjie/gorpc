package fs

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

// Move move `src` to `dest`
//
// the behavior of fs.Move is consistent with bash shell `mv` command:
//
// when move a file, actions are following:
// ------------------------------------------------------------------------------------------------
// | No. | src existed | src type | dst existed | dst type | behavior                             |
// ------------------------------------------------------------------------------------------------
// | 1   | False       | -        | -           | -        | error: No such file or directory     |
// ------------------------------------------------------------------------------------------------
// | 2   | True        | File     | False       | -        | if dir(dst) existed:                 |
// |     |             |          |             |          | - Yes, is dir, mv `src` to dir(dst)  |
// |     |             |          |             |          | - Yes, not dir, err: Not a directory |
// |     |             |          |             |          | - No, err: No such file or directory |
// ------------------------------------------------------------------------------------------------
// | 3   | True        | File     | True        | Folder   | if dst/basename(src) existed:        |
// |     |             |          |             |          | - Yes, mv `src` to dst/basename(src) |
// |     |             |          |             |          | - No, mv `src` to dst/basename(src)  |
// ------------------------------------------------------------------------------------------------
//
// when move a directory, actions are following:
// ------------------------------------------------------------------------------------------------
// | 4   | True        | File     | True        | File     | mv `src` to dst                      |
// ------------------------------------------------------------------------------------------------
// | 5   | True        | Folder   | False       | -        | if dir(dst) existed:                 |
// |     |             |          |             |          | - Yes, is dir, mv `src` to dir(dst)  |
// |     |             |          |             |          | - Yes, not dir, err: File Exists     |
// |     |             |          |             |          | - No, err: No such file or directory |
// ------------------------------------------------------------------------------------------------
// | 6   | True        | Folder   | True        | File     | error: File Already Existed          |
// ------------------------------------------------------------------------------------------------
// | 7   | True        | Folder   | True        | Folder   | t = dst/basename(src), if t existed: |
// |     |             |          |             |          | - Yes, t empty, mv src to t          |
// |     |             |          |             |          | -      t notempty, err: t Not empty  |
// |     |             |          |             |          | - No, mv src to t                    |
// ------------------------------------------------------------------------------------------------
//
// Why keep the behavior consistent? It makes the usage much more friendly when it behaves as users expected.
func Move(src, dst string) error {

	var (
		inf os.FileInfo
		err error
	)

	// check whether `src` is valid or not
	if inf, err = os.Lstat(src); err != nil {
		return err
	}

	// move directory
	if inf.IsDir() {
		return moveDirectory(src, dst)
	}

	// move file
	return moveFile(src, dst)
}

// moveFile move a file `src` to `dst`
//
// `src` is a normal file, dst can be a file or directory.
// 1. if `dst` not existed
// 	- if dir(dst) existed and is a directory, then move `src` under dir(dst),
//  - if dir(dst) existed and not a directory, return err: &PathError("lstat", dir(dst), syscall.EEXIST}
//  - if dir(dst) not existed, return err: &PathError("lstat", dir(dst), os.ENOENT}
// 2. if `dst` existed
// - if dst is a normal file, rename src to dst
// - if dst is a folder, rename src to dst/basename(src)
func moveFile(src, dst string) error {

	dstInf, err := os.Lstat(dst)

	// if dst existed

	if err == nil {
		if !dstInf.IsDir() {
			return os.Rename(src, dst)
		}

		p := path.Join(dst, filepath.Base(src))
		return os.Rename(src, p)
	}

	if !os.IsNotExist(err) {
		return err
	}

	// if dst not existed

	p := filepath.Dir(dst)
	if inf, err := os.Lstat(p); err != nil {
		return err
	} else {
		if !inf.IsDir() {
			return &os.PathError{"lstat", p, syscall.EEXIST}
		}
		return os.Rename(src, dst)
	}
}

// moveDirectory move a directory `src` to `dst`
//
// `src` is a directory, dst should always be a directory.
// 1. if `dst` existed
// 	- if `dst` is not a directory, return error &PathError{"lstat", dst, os.EEXIST}
// 	- if `dst` is a directory
//		- if dst/basename(src) is empty, then rename src to dst/basename(src)
//		- if dst/basename(src) not empty, return error &PathError{"lstat", dst, syscall.ENOTEMPTY}
// 2. if `dst` not existed
// 	- if dir(dst) existed, rename src to dst
// 	- if dir(dst) not existed, return error &PathError{"lstat", dst, syscall.ENOENT}
func moveDirectory(src, dst string) error {

	dstInf, err := os.Lstat(dst)

	// if dst existed
	if err == nil {

		if !dstInf.IsDir() {
			return &os.PathError{"lstat", dst, syscall.EEXIST}
		}

		target := path.Join(dst, filepath.Base(src))
		inf, err := os.Lstat(target)
		if err != nil {
			if os.IsNotExist(err) {
				return os.Rename(src, target)
			}
			return err
		}

		if !inf.IsDir() {
			return &os.PathError{"lstat", target, syscall.EEXIST}
		}

		files, err := ioutil.ReadDir(target)
		if err != nil {
			return err
		}
		if len(files) != 0 {
			return &os.PathError{"lstat", target, syscall.ENOTEMPTY}
		}

		if err = os.RemoveAll(target); err != nil {
			return err
		}
		return os.Rename(src, target)
	}

	// if dst not existed
	inf, err := os.Lstat(filepath.Dir(dst))
	if err != nil {
		return err
	}
	if !inf.IsDir() {
		return &os.PathError{"lstat", filepath.Base(dst), syscall.EEXIST}
	}
	return os.Rename(src, dst)
}
