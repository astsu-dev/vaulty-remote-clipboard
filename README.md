# Clipsync - sync your clipboard between PC and Mobile

## Build

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
make build_macos
```

After the build is complete, you can find the `clipsync` executable file in the `bin` directory.

## How to run

### Generate SSL certificate and private key

```bash
make gen_cert
```

### Create config file

Create config file by the following path `~/.config/clipsync/config.toml` if you use Linux or MacOS, or `./config.toml` if you use Windows.

Paste the following content into the config file:

```yaml
apiKeyForGetClipboard = "api_key_for_client"
apiKeyForSetClipboard = "api_key_for_server"
httpsCertificateFilePath = "cert.pem" # path to the generated certificate from the previous step
httpsCertificateKeyFilePath = "key.pem" # path to the generated private key from the previous step
port = 8090
```

### Run the server:

```bash
./bin/clipsync
```

If all is well, you will see the IP address and port on which the server is running.
You should enter this data in the mobile application.

### Non-default config location

You can also specify the path to the config file if it is not in the default location:

```bash
./bin/clipsync /path/to/config.toml
```
