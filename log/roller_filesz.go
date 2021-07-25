package log

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	b = 1 << (10 * iota)
	kb
	mb
	gb
)

// only Byte/KiloByte/MegaByte/GigaByte supported, other unit will be treated as Byte
var sizePattern = regexp.MustCompile(`^(\d+)([kKmMgG]?[bB]?)$`)

// parse string 'size' and return number in unit bytes
func filesize(size string) (bytes int, err error) {
	if !sizePattern.MatchString(size) {
		return 0, fmt.Errorf("sizePattern mot match")
	}

	s := sizePattern.FindStringSubmatch(size)
	if len(s) != 3 {
		return 0, fmt.Errorf("capture groups invalid")
	}

	sz, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, err
	}

	unit := strings.ToLower(s[2])
	switch unit {
	case "b":
		sz *= b
	case "k", "kb":
		sz *= kb
	case "m", "mb":
		sz *= mb
	case "g", "gb":
		sz *= gb
	default:
		sz *= b
	}

	return sz, nil
}

type rollerByFileSZ struct {
	logger *logger
}

func (r *rollerByFileSZ) Roll(oldf Writer) (newf Writer, err error) {
	return
}
