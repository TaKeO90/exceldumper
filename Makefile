BINARY=exceldumper
SERVER=srv
SERVERPATH=./server/main.go

.PHONY:build
build:
	go mod download
	go build .

.PHONY:server
server:
	go mod download
	go build -o $(SERVER) $(SERVERPATH)

.PHONY:cleansrv
cleansrv: 
	rm $(SERVER)

.PHONY:clean
clean:
	rm $(BINARY)
