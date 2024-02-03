package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
)

type ResourceErr struct {
	Resource string
	Code     int
}

func (m ResourceErr) Error() string {
	return fmt.Sprintf("resource %s: %d", m.Resource, m.Code)
}

func (m ResourceErr) Is(target error) bool {
	if other, ok := target.(ResourceErr); ok {
		return m.Resource == other.Resource && m.Code == other.Code
	}
	return false
}

func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("fileChecker: %w", err)
	}
	f.Close()
	return nil
}

func main() {
	// ex 1
	data := []byte("not a zip file")
	file := bytes.NewReader(data)

	_, err := zip.NewReader(file, int64(len(data)))
	if errors.Is(err, zip.ErrFormat) {
		fmt.Println("ErrFormat")
	}

	// ex 2 - unwrap error. However, errors.Is and errors.As are more idiomatic
	err = fileChecker("non-existing-file.txt")
	if err != nil {
		fmt.Println(err)
		if wrappedErr := errors.Unwrap(err); wrappedErr != nil {
			fmt.Println(wrappedErr)
		}
	}

	// ex 3 - errors.Is - checks if any error in the chain of wrapped errors matches the target
	err = fileChecker("non-existing-file.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
	}

	// ex 4 - custom Is() method
	err = ResourceErr{"db", 1}
	if errors.Is(err, ResourceErr{"db", 1}) {
		fmt.Println("db error")
	}

	// ex 5 - errors.As - checks whether the error has a specific type
	err = ResourceErr{"db", 1}
	var rErr ResourceErr
	if errors.As(err, &rErr) {
		fmt.Println(rErr)
	}

	// panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from panic:", r)
			}
		}()
		panic("unexpected error occurred")
	}()

	fmt.Println("End of error handling")
}
