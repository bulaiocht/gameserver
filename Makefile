client:
	@go build -o .\bin\client.exe .\game_client\main.go
	@bin\client

server:
	@go build -o .\bin\server.exe .\game_server\main.go
	@bin\server
