package nasin

import "os"
import "fmt"
// TODO: possibly fork the official plugin module and add support for other
// operating systems? perhaps enhance the Lookup function with
// the generic extract function we have here for extra type safety goodness.
import "plugin"
import "strings"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo"

type expectsFunc     func () (int, int, int)
type nameFunc        func () string
type descriptionFunc func () string
type backendFactory  func () (tomo.Backend, error)
type themeFactory    func () tomo.Theme

var factories []backendFactory
var theme tomo.Theme

func loadPlugins () {
	// TODO: do not hardcode all these paths here, have separate files that
	// build on different platforms that set these paths.
	
	pathVariable := os.Getenv("NASIN_PLUGIN_PATH")
	paths := strings.Split(pathVariable, ":")
	paths = append (
		paths,
		"/usr/lib/nasin/plugins",
		"/usr/local/lib/nasin/plugins")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		paths = append (
			paths,
			filepath.Join(homeDir, ".local/lib/nasin/plugins"))
	}

	for _, dir := range paths {
		loadPluginsFrom(dir)
	}
}

func loadPluginsFrom (dir string) {
	entries, err := os.ReadDir(dir)
	// its no big deal if one of the dirs doesn't exist
	if err != nil { return }

	for _, entry := range entries {
		if entry.IsDir() { continue }
		if filepath.Ext(entry.Name()) != ".so" { continue }
		pluginPath := filepath.Join(dir, entry.Name())
		loadPlugin(pluginPath)
	}
}

func loadPlugin (path string) {
	die := func (reason string) {
		println (
			"nasin: could not load plugin at ",
			path + ":", reason)
	}

	plugin, err := plugin.Open(path)
	if err != nil {
		die(err.Error())
		return
	}

	// check for and obtain basic plugin functions
	expects, ok := extract[expectsFunc](plugin, "Expects")
	if !ok { die("does not implement Expects() (int, int, int)"); return }
	name, ok := extract[nameFunc](plugin, "Name")
	if !ok { die("does not implement Name() string"); return }
	_, ok = extract[descriptionFunc](plugin, "Description")
	if !ok { die("does not implement Description() string"); return }

	// check for version compatibility
	// TODO: have exported version type in tomo base package, and have a
	// function within that that gives the current tomo/nasin version. call
	// that here.
	major, minor, patch := expects()
	currentVersion := version { 0, 0, 0 }
	pluginVersion  := version { major, minor, patch }
	if !pluginVersion.CompatibleABI(currentVersion) {
		die (
			"plugin (" + pluginVersion.String() +
			") incompatible with nasin/tomo version (" +
			currentVersion.String() + ")")
		return
	}

	// if it's a backend plugin...
	newBackend, ok := extract[backendFactory](plugin, "NewBackend")
	if ok { factories = append(factories, newBackend) }

	// if it's a theme plugin...
	newTheme, ok := extract[themeFactory](plugin, "NewTheme")
	if ok { theme = newTheme() }

	println("nasin: loaded plugin", name())
}

func extract[T any] (plugin *plugin.Plugin, name string) (value T, ok bool) {
	symbol, err := plugin.Lookup(name)
	if err != nil { return }
	value, ok = symbol.(T)
	return
}

type version [3]int

func (version version) CompatibleABI (other version) bool {
	return version[0] == other[0] && version[1] == other[1]
}

func (version version) String () string {
	return fmt.Sprint(version[0], ".", version[1], ".", version[2])
}
