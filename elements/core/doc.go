// Package core provides tools that allow elements to easily fulfill common
// interfaces without having to duplicate a ton of code. Each "core" is a type
// that can be embedded into an element directly, working to fulfill a
// particular interface. Each one comes with a corresponding core control, which
// provides an interface for elements to exert control over the core. Core
// controls should be kept private.
package core
