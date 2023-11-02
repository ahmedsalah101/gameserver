client:
	@go build -o bin/client game_client/main.go
	@./bin/client
server:
	@go build -o bin/server game_server/main.go
	@./bin/server
compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
live:
	gin -p 3000 run main.go
