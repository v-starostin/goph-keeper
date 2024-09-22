.PHONY: run
run:
	export DB_CONN_STRING=.db_conn && export SECRET=.secret && export ENCRYPTION_KEY=.secret && go run cmd/main.go

.PHONY: gen-proto
gen-proto:
	protoc \
		--go_out=. --go_opt=module=github.com/v-starostin/goph-keeper \
  		--go-grpc_out=. --go-grpc_opt=module=github.com/v-starostin/goph-keeper \
  		api/api.proto

.PHONY: run-db
run-db:
	@docker run \
		-d \
		-v `pwd`/db/migration:/docker-entrypoint-initdb.d/ \
		--rm \
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:16
