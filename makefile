
all: compile run

compile:
	go build server.go

run:
	./server 8080