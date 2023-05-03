package nasin

import "os"
// TODO: possibly fork the official plugin module and add support for other
// operating systems? perhaps enhance the Lookup function with
// the generic extract function we have here for extra type safety goodness.
import "plugin"
import "path/filepath"
import "tomo"

type backendFactory  func () (tomo.Backend, error)
var factories []backendFactory
var theme tomo.Theme

var pluginPaths []string

func loadPlugins () {
	for _, dir := range pluginPaths {
		if dir != "" {
			loadPluginsFrom(dir)
		}
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
		println("nasin: could not load plugin at", path + ":", reason)
	}

	plugin, err := plugin.Open(path)
	if err != nil {
		die(err.Error())
		return
	}

	// check for and obtain basic plugin functions
	expects, ok := extract[func () tomo.Version](plugin, "Expects")
	if !ok { die("does not implement Expects() tomo.Version"); return }
	name, ok := extract[func () string](plugin, "Name")
	if !ok { die("does not implement Name() string"); return }
	_, ok = extract[func () string](plugin, "Description")
	if !ok { die("does not implement Description() string"); return }

	// check for version compatibility
	pluginVersion  := expects()
	currentVersion := tomo.CurrentVersion()
	if !pluginVersion.CompatibleABI(currentVersion) {
		die (
			"plugin (" + pluginVersion.String() +
			") incompatible with tomo/nasin version (" +
			currentVersion.String() + ")")
		return
	}

	// if it's a backend plugin...
	newBackend, ok := extract[func () (tomo.Backend, error)](plugin, "NewBackend")
	if ok { factories = append(factories, newBackend) }

	// if it's a theme plugin...
	newTheme, ok := extract[func () tomo.Theme](plugin, "NewTheme")
	if ok { theme = newTheme() }

	println("nasin: loaded plugin", name())
}

func extract[T any] (plugin *plugin.Plugin, name string) (value T, ok bool) {
	symbol, err := plugin.Lookup(name)
	if err != nil { return }
	value, ok = symbol.(T)
	return
}
