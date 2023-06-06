build:
	@go build -o bin/cashe

run: build
	@./bin/cashe

runfollower: build
	@./bin/cashe --listenaddr :4000 --leaderaddr :3000

ct: 
	@go build -o bin/client client/main.go
	@./bin/client

test:
	@go test -v ./...