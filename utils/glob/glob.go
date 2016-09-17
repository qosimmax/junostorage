package glob

import (
	glob "path/filepath"
)

func Match(pattern, name string) (matched bool, err error) {
	matched, err = glob.Match(pattern, name)
	return
}
