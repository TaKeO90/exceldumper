BINARY=exceldumper


.PHONY:build
build:
	go mod download
	go build .


.PHONY:clean
clean: 
	rm exceldumper
