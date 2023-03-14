start-dependencies:
	docker-compose up -d 

start-api:
	go run cmd/api/main.go