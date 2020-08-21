.PHONY: all
all: temperature pi_copy

temperature: main.go
	env GOOS=linux GOARCH=arm GOARM=6 go build

.PHONY: pi_copy
pi_copy: main.go
	scp temperature pi@192.168.1.105:
