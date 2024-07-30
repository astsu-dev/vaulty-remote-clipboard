fmt:
	go fmt ./...
vet:
	go vet ./...
lint: vet
	golangci-lint run
run:
	go run main.go

# GUI build
build_darwin:
	fyne-cross darwin -arch=arm64
build_windows:
	fyne-cross windows -arch=amd64
build_linux:
	fyne-cross linux -arch=amd64

# CLI build
build_darwin_cli:
	GOOS=darwin GOARCH=arm64 go build cmd/remclip_cli/remclip_cli.go -o bin/remclip
build_windows_cli:
	GOOS=windows GOARCH=amd64 go build cmd/remclip_cli/remclip_cli.go -o bin/remclip.exe
build_linux_cli:
	GOOS=linux GOARCH=amd64 go build cmd/remclip_cli/remclip_cli.go -o bin/remclip
