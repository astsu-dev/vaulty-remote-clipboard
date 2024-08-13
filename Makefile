fmt:
	go fmt ./...
vet:
	go vet ./...
lint: vet
	golangci-lint run
run:
	go run main.go
test:
	go test -v ./...

# GUI build
build_clean:
	rm -rf fyne-cross && rm -rf dist && mkdir dist
build_darwin: build_clean
	fyne-cross darwin -arch=arm64 \
	&& mkdir fyne-cross/dist/darwin-arm64/dist \
	&& mv "fyne-cross/dist/darwin-arm64/Vaulty Remote Clipboard.app" fyne-cross/dist/darwin-arm64/dist \
	&& ln -s /Applications fyne-cross/dist/darwin-arm64/dist/Applications \
	&& hdiutil create -volname "Vaulty Remote Clipboard" -srcfolder "fyne-cross/dist/darwin-arm64/dist" -ov -format UDZO "./dist/Vaulty Remote Clipboard ARM.dmg"
build_darwin_x86: build_clean
	fyne-cross darwin -arch=amd64 \
	&& mkdir fyne-cross/dist/darwin-amd64/dist \
	&& mv "fyne-cross/dist/darwin-amd64/Vaulty Remote Clipboard.app" fyne-cross/dist/darwin-amd64/dist \
	&& ln -s /Applications fyne-cross/dist/darwin-amd64/dist/Applications \
	&& hdiutil create -volname "Vaulty Remote Clipboard" -srcfolder "fyne-cross/dist/darwin-amd64/dist" -ov -format UDZO "./dist/Vaulty Remote Clipboard x86.dmg"
build_windows: build_clean
	fyne-cross windows -arch=amd64 \
	&& docker run --rm -i -v $$PWD:/work amake/innosetup innosetup.iss \
	&& mv ./Output/* dist/ && rm -rf ./Output
build_linux: build_clean
	fyne-cross linux -arch=amd64 \
	&& mv fyne-cross/dist/linux-amd64/* ./dist

# CLI build
build_darwin_cli:
	GOOS=darwin GOARCH=arm64 go build cmd/remclip_cli/remclip_cli.go -o bin/remclip
build_windows_cli:
	GOOS=windows GOARCH=amd64 go build cmd/remclip_cli/remclip_cli.go -o bin/remclip.exe
build_linux_cli:
	GOOS=linux GOARCH=amd64 go build cmd/remclip_cli/remclip_cli.go -o bin/remclip
