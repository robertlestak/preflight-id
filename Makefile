bin: bin/preflight-id_darwin_amd64 bin/preflight-id_linux_amd64 bin/preflight-id_windows_amd64.exe
bin: bin/preflight-id_darwin_arm64 bin/preflight-id_linux_arm64 bin/preflight-id_windows_arm64.exe

bin/preflight-id_darwin_amd64:
	@mkdir -p bin
	@echo "Compiling preflight-id..."
	GOOS=darwin GOARCH=amd64 go build -o $@ cmd/preflight-id/*.go

bin/preflight-id_darwin_arm64:
	@mkdir -p bin
	@echo "Compiling preflight-id..."
	GOOS=darwin GOARCH=arm64 go build -o $@ cmd/preflight-id/*.go

bin/preflight-id_linux_amd64:
	@mkdir -p bin
	@echo "Compiling preflight-id..."
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/preflight-id/*.go

bin/preflight-id_linux_arm64:
	@mkdir -p bin
	@echo "Compiling preflight-id..."
	GOOS=linux GOARCH=arm64 go build -o $@ cmd/preflight-id/*.go

bin/preflight-id_windows_amd64.exe:
	@mkdir -p bin
	@echo "Compiling preflight-id..."
	GOOS=windows GOARCH=amd64 go build -o $@ cmd/preflight-id/*.go

bin/preflight-id_windows_arm64.exe:
	@mkdir -p bin
	@echo "Compiling preflight-id..."
	GOOS=windows GOARCH=arm64 go build -o $@ cmd/preflight-id/*.go

.PHONY: install
install: bin
	@echo "Installing preflight-id..."
	@scp bin/preflight-id_$$(go env GOOS)_$$(go env GOARCH) /usr/local/bin/preflight-id