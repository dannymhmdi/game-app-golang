first implement this part in logger package
```go
package logger

import (
	"os"
	"log"
)
func SetUpFile(filename string) (*os.File, error) {
	file, oErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if oErr != nil {
		return nil, oErr
	}
	log.SetOutput(file)
	return file, nil
}
```

second implement this part in main function:

```go
package main

import (
	"mymodule/logger"
	"log"
)

logFile, sErr := logger.SetUpFile("errors.log")
if sErr != nil {
	log.Fatal("failed to setup logger file")
}
defer logFile.Close()
config.Load()
```