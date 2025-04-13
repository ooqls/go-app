package app

import "os"

func mustReadFile(p string) []byte {
	b, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}

	return b
}
