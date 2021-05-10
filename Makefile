install: 
	go install -i

build: 
	go build -o leetgo.exe

get: 
	go build -o leetgo.exe
	leetgo.exe get 10
	

