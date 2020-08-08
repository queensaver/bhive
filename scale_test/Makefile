scale_test: main.go
	env GOOS=linux GOARCH=arm GOARM=5 go build

.PHONY: pi_copy
pi_copy: main.go
	scp scale_test pi@192.168.1.105:

