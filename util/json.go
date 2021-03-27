package util

import (
	"encoding/json"
	"os"
)

func LoadJSON(filename string, target interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(target)

	return
}

func SaveJSON(filename string, source interface{}) (err error) {
	// Yes, this is a very naive implementation (no atomic overwrites etc.),
	// but for our purposes we don't care if a file gets corrupted
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(source)

	return
}
