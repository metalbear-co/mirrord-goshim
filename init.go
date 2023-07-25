package mirrord_goshim

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ebitengine/purego"
)

func init() {
	envName := getEnvName()
	if envName == "" {
		fmt.Fprintf(os.Stderr, "Failed to inject mirrord layer: GOOS=%s is not supported.\n", runtime.GOOS)
		return
	}

	layerPath := getLayerPath(envName)
	if layerPath == "" {
		fmt.Fprintf(os.Stderr, "Failed to inject mirrord layer: layer file not found in %s\n", envName)
		return
	}

	_, err := purego.Dlopen(layerPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to inject mirrord layer %s: %e\n", layerPath, err)
		return
	}

	fmt.Fprintf(os.Stderr, "Injected mirrord layer %s\n", layerPath)
}

func getEnvName() string {
	switch runtime.GOOS {
	case "darwin":
		return "DYLD_INSERT_LIBRARIES"
	case "linux":
		return "LD_PRELOAD"
	default:
		return ""
	}
}

func getLayerPath(envName string) string {
	env := os.Getenv(envName)
	paths := strings.Split(env, ":")

	for i := len(paths) - 1; i >= 0; i-- {
		if strings.Contains(paths[i], "libmirrord_layer") {
			return paths[i]
		}
	}

	return ""
}
