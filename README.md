# Vaulty Remote Clipboard - copy your passwords from the Vaulty mobile app to your computer

The main philosophy of the [Vaulty](https://github.com/astsu-dev/vaulty-mobile) mobile application is to store your passwords only on your smartphone from the security purposes.
This is why the application does not have a web version or a desktop version.
However, sometimes you need to copy your password from the Vaulty to your computer.
This is where the Vaulty Remote Clipboard application comes in handy.
It allows you to securely copy your passwords from the Vaulty mobile application to your computer using the local Wi-Fi network.

## GUI

### Installation

Download the latest version of the application from the [Releases](https://github.com/astsu-dev/vaulty-remote-clipboard/releases) page.

### How to use the application?

Run the downloaded executable file. The application will prompt the password from you. Enter the same as in the mobile application.
By default the server is running on the `8090` port. You can change this value in the port field. If you change the port, don't forget to click the "Save" button.
To run the server, click the "Start" button. If all is well, you will see the log message indicates the server is running below the start button.
If you want to stop the server, click the "Stop" button. You can close the application and it will be minimized to the system tray.

### How it works?

The application starts the UDP server on the specified port.
The entered password is used to encrypt and decrypt the data between the mobile application and the server using AES GCM encryption algorithm.
This is why you need to enter the same password on both devices.
The mobile application sends the encrypted UDP datagram with the desired data to copy to the server when you click the "Send to remote" button (computer icon) in your app.
It is not needed for the mobile application to know the exact IP address of the server as it uses the broadcast address of the local network.
It means that you must be connected to the same local network to use the remote clipboard feature.
As the broadcast address is used, the encrypted datagram is sent to all devices in the local network.
However only devices with the same encryption password can decrypt the datagram.
You should use the strong password to prevent the data interception by other devices in the local network.
After the server receives the datagram, it decrypts it and set the decrypted data to the system clipboard.

### Build

If you want to build the GUI application by yourself and not use the pre-built binary
from the [Releases](https://github.com/astsu-dev/vaulty-remote-clipboard/releases) page you can use the following commands for your OS.

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

## CLI

There is also an alternative way to use the application. You can use the CLI version of the application.

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

### Build

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
