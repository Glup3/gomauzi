run:
	go run main.go


dev:
	air


db:
	docker-compose up -d


db_down:
	docker-compose down


db_reset:
	docker-compose down
	docker-compose up -d


test:
	go test -v ./...