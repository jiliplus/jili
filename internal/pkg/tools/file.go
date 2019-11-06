package tools

import "os"

// IsExist return true if path is exist, otherwise return false
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Create writes bytes to new file filename
func Create(filename string, bytes []byte) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		panic(err)
	}
}
