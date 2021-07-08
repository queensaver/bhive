bhive: main.go
	go build
	
bhive_arm64: main.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o bhive_arm64 -ldflags="-s -w" -v
