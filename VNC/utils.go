package VNC

import (
	"fmt"
)

func failFatal(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
