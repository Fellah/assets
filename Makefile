install: format
	go vet
	go test
	go install -gcflags "-N -l"

format:
	go fmt
