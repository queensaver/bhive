scale: main.go
	env GOOS=linux GOARCH=arm GOARM=5 go build -tags=gpiomem

.PHONY: pi_copy
pi_copy: main.go
	scp scale pi@192.168.1.105:

