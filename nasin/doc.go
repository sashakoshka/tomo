// Package nasin provides a higher-level framework for Tomo applications. Nasin
// also automatically handles themes, backend instantiation, and plugin support.
//
// Backends and themes are typically loaded through plugins. For now, plugins
// are only supported on UNIX-like systems, but this will change in the future.
// Nasin will attempt to load all ".so" files in these directories as plugins:
//
//   - /usr/lib/nasin/plugins
//   - /usr/local/lib/nasin/plugins
//   - $HOME/.local/lib/nasin/plugins
//
// It will also attempt to load all ".so" files in the directory specified by
// the NASIN_PLUGIN_PATH environment variable.
//
// Plugins must export the following functions at minimum:
//
//   + Expects() tomo.Version
//   + Name() string
//   + Description() string
//
// Expects() must return the version of Tomo/Nasin it was built for. Nasin will
// automatically figure out if the plugin has a compatible ABI with the current
// version and refuse to load it if not. Name() and Description() return a short
// plugin name and a description of what a plugin does, respectively. Plugins
// must not attempt to interact with Tomo/Nasin within their init functions.
//
// If a plugin provides a backend, it must export this function:
//
//   NewBackend() (tomo.Backend, error)
//
// This function must attempt to initialize the backend, and return it if
// successful. Otherwise, it should clean up all resources and return an error
// explaining what caused the backend to fail to initialize. The first backend
// that does not throw an error will be used.
//
// If a plugin provides a theme, it must export this function:
//
//   NewTheme() tomo.Theme
//
// This just creates a new theme and returns it.
//
// For information on how to create plugins with Go, visit:
// https://pkg.go.dev/plugin
package nasin
