//go:build linux || darwin || freebsd

package nasin

import "os"
import "strings"
import "path/filepath"

func init () {
	pathVariable := os.Getenv("NASIN_PLUGIN_PATH")
	pluginPaths = strings.Split(pathVariable, ":")
	pluginPaths = append (
		pluginPaths,
		"/usr/lib/nasin/plugins",
		"/usr/local/lib/nasin/plugins")
	homeDir, err := os.UserHomeDir()
	if err == nil {
		pluginPaths = append (
			pluginPaths,
			filepath.Join(homeDir, ".local/lib/nasin/plugins"))
	}
}
