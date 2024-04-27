migrate:
	@migrate -path ./migrations/ -database "postgres://postgres:coffice@213.226.127.170/mcatalogue?sslmode=disable" up

down:
	@migrate -path ./migrations/ -database "postgres://postgres:coffice@213.226.127.170/mcatalogue?sslmode=disable" down

