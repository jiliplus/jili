package tools

import "os"

// IsExist return true if path is exist, otherwise return false
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
