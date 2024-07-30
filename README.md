# Vaulty Remote Clipboard - copy your passwords from the Vaulty mobile app to your computer

## GUI

### Build

> You can skip this step if you use the pre-built binary from the [Releases](https://github.com/astsu-dev/vaulty-remote-clipboard/releases) page.

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

### How to use the application

Run the built executable file. The application will prompt the password from you. Enter the same as in the mobile application.
By default the server is running on the `8090` port. You can change this value in the port field. If you change the port, don't forget to click the "Save" button.
To run the server, click the "Start" button. If all is well, you will see the log message indicates the server is running below the start button.
If you want to stop the server, click the "Stop" button.

## CLI

### Build

> You can skip this step if you use the pre-built binary from the [Releases](https://github.com/astsu-dev/vaulty-remote-clipboard/releases) page.

Windows:

```bash
make build_windows_cli
```

Linux:

```bash
make build_linux_cli
```

MacOS:

```bash
make build_darwin_cli
```

After the build is complete, you can find the executable file in the `./bin` directory.

### Create config file

Create config file by the following path `~/.config/remclip/config.toml` if you use Linux or MacOS, or `./config.toml` if you use Windows.

Paste the following content into the config file:

```yaml
port = 8090
serverType = udp
```

### Run the server:

```bash
./bin/remclip
```

The application will prompt the password from you. Enter the same as in the mobile application.
If all is well, you will see the log message indicates the server is running.

### Non-default config location

You can also specify the path to the config file if it is not in the default location:

```bash
./bin/remclip /path/to/config.toml
```
