gen-pb:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative v1/otabe.proto

migrateup:
	 migrate -path db/migration -database "mysql://root:Hannamysql.1518@tcp(localhost:49530)/otabe" -verbose up

migratedown:
	 migrate -path db/migration -database "mysql://root:Hannamysql.1518@tcp(localhost:49530)/otabe" -verbose down

# Create 2 file up/down for migration:
# migrate create -ext sql -dir db/migration -seq init_schema
