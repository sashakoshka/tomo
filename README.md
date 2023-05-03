# ![tomo](assets/screenshot.png)

This repository is [mirrored on GitHub](https://github.com/sashakoshka/tomo).

Please note: Tomo is in early development. Some features may not work properly,
and its API may change without notice.

Tomo is a GUI toolkit written in pure Go with minimal external dependencies. It
makes use of Go's unique language features to do more with less.

Nasin is an application framework that runs on top of Tomo. It supports plugins
which can extend any application with backends, themes, etc.

## Usage

Before you start using Tomo, you need to install a backend plugin. Currently,
there is only an X backend. You can run ./scripts/install-backends.sh to install
it. It will be placed in `~/.local/lib/nasin/plugins`.

You can find out more about how to use Tomo and Nasin by visiting the examples
directory, or pull up the documentation by running `godoc` within the
repository. You can also view it on the web on
[pkg.go.dev](https://pkg.go.dev/git.tebibyte.media/sashakoshka/tomo) (although
it may be slightly out of date).
