<img align="left" width="80" height="80" src="assets/icon.png" alt="App icon" />

# Build

If you want to build the application by yourself and not use the pre-built binary
from the [Releases](https://github.com/astsu-dev/vaulty-remote-clipboard/releases/latest) page you can use the following commands for your OS.

Before you start the build process, you need to make sure that you have installed all the necessary dependencies:

- go 1.22
- make
- fyne-cross tool. How to install it you can find [here](https://docs.fyne.io/started/cross-compiling.html). fyne-cross requires the docker to be installed on your machine.

Windows:

```bash
make build_windows
```

Linux:

```bash
make build_linux
```

MacOS:

```bash
make build_darwin
```

After the build is complete, you can find the executable file in the `./fyne-cross/bin/<os-name>` directory.
