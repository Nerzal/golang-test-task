start-dependencies:
	docker-compose up -d 

start-api:
	go run cmd/api/main.go

start-processor:
	go run cmd/processor/main.go

start-reporting-api:
	go run cmd/reporting-api/main.go