package gmx

// pkg/os instrumentation

import (
	"os"
)

func init() {
	Registry("os")("args", osArgs)
}

func osArgs() interface{} {
	return os.Args
}
