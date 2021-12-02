gen:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb

migrateup:
	 migrate -path db/migration -database "mysql://root:Hannamysql.1518@tcp(localhost:49530)/otabe" -verbose up

migratedown:
	 migrate -path db/migration -database "mysql://root:Hannamysql.1518@tcp(localhost:49530)/otabe" -verbose down

# Create 2 file up/down for migration:
# migrate create -ext sql -dir db/migration -seq init_schema
