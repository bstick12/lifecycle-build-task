package task

import (
	"fmt"
	"os"
	"strings"
)

func ResetEnv(origEnv []string) {
	os.Clearenv()
	for _, pair := range origEnv {
		i := strings.Index(pair[1:], "=") + 1
		if err := os.Setenv(pair[:i], pair[i+1:]); err != nil {
			panic(fmt.Sprintf("Setenv(%q, %q) failed during reset: %v", pair[:i], pair[i+2:], err))
		}
	}
}
