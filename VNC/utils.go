package VNC

import (
	"errors"
	"fmt"
	"log"
)

func failFatal(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
func logError(format string, args ...interface{}) error {
	log.Printf(format, args)
	return errors.New(fmt.Sprintf(format, args))
}
