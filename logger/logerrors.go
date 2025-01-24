package logger

import (
	"log"
	"os"
)

func SetUpFile(filename string) (*os.File, error) {
	file, oErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if oErr != nil {
		return nil, oErr
	}
	log.SetOutput(file)
	return file, nil
}
