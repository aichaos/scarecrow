package scarecrow

import "os"

// MakeDirectory makes a directory if it doesn't already exist.
func MakeDirectory(path string) {
	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, 0755)
	}
}
