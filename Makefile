api:
	go run cmd/api/notify-api/main.go

migrateup:
	go run cmd/db/migrate_up/main.go

migratedown:
	go run cmd/db/migrate_down/main.go

migrateforce:
	go run cmd/db/migrate_force/main.go

dispatcher:
	go run cmd/consumer/dispatcher/main.go
