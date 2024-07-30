fmt:
	go fmt ./...
vet:
	go vet ./...
lint: vet
	golangci-lint run
build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/remclip.exe cmd/rempclip/rempclip.go
build_macos:
	GOOS=darwin GOARCH=arm64 go build -o bin/remclip cmd/rempclip/rempclip.go
build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/rempclip cmd/rempclip/rempclip.go
build_gen_cert_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/generate_cert.exe cmd/generate_cert/generate_cert.go
build_gen_cert_macos:
	GOOS=darwin GOARCH=arm64 go build -o bin/generate_cert cmd/generate_cert/generate_cert.go
build_gen_cert_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/generate_cert cmd/generate_cert/generate_cert.go
build_bundle: build_macos build_windows
	cp bin/rempclip bin/rempclip.exe config.toml rempclip/ && zip -r rempclip.zip rempclip
run:
	go run cmd/rempclip/rempclip.go config.toml
gen_cert:
	go run cmd/generate_cert/generate_cert.go --host example.com
