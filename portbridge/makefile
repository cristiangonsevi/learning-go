
release-all: release-linux-x86 release-linux-arm64 release-windows-x86 release-windows-arm64 release-darwin-arm64

release-linux-x86:
	GOOS=linux GOARCH=amd64 go build -o ./releases/pb-linux-x86 ./cmd/portbridge
release-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o ./releases/pb-linux-arm64 ./cmd/portbridge
release-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -o ./releases/pb-windows-arm64.exe ./cmd/portbridge
release-windows-x86:
	GOOS=windows GOARCH=amd64 go build -o ./releases/pb-windows-x86.exe ./cmd/portbridge
release-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./releases/pb-darwin-arm64 ./cmd/portbridge
