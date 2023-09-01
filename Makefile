build-linux:
	GOOS=linux GOARCH=amd64 go build  -o build/kong-keyless main.go
build-win:
	GOOS=windows GOARCH=amd64 go build -o build/kong-keyless.exe main.go
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o build/kong-keyless.exe main.go
build-mac-silicon:
	GOOS=darwin GOARCH=arm64 go build -o build/kong-keyless.exe main.go
build-native:
	go build -o build/kong-keyless.exe main.go
