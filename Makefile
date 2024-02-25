server:
	@go build -o .\bin\server.exe .\game_server\server.go
	@bin\server

client:
	@go build -o .\bin\client.exe .\game_client\client.go
	@bin\client
