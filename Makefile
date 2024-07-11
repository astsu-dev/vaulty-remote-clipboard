fmt:
	go fmt ./...
vet:
	go vet ./...
lint: vet
	golangci-lint run
build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/clipsync.exe cmd/clipsync/clipsync.go
build_macos:
	GOOS=darwin GOARCH=arm64 go build -o bin/clipsync cmd/clipsync/clipsync.go
build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/clipsync cmd/clipsync/clipsync.go
build_gen_cert_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/generate_cert.exe cmd/generate_cert/generate_cert.go
build_gen_cert_macos:
	GOOS=darwin GOARCH=arm64 go build -o bin/generate_cert cmd/generate_cert/generate_cert.go
build_gen_cert_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/generate_cert cmd/generate_cert/generate_cert.go
build_bundle: build_macos build_windows
	cp bin/clipsync bin/clipsync.exe config.toml clipsync/ && zip -r clipsync.zip clipsync
run:
	go run cmd/clipsync/clipsync.go config.toml
gen_cert:
	go run cmd/generate_cert/generate_cert.go --host example.com
