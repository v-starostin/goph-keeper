.PHONY: run
run:
	export DB_CONN_STRING=.db_conn && \
	export SECRET=.secret && \
 	export ENCRYPTION_KEY=.encryption_key && \
 	go run cmd/main.go

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
		-p 5432:5432 \
		--name db \
		-e POSTGRES_DB=postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		postgres:16

.PHONY: run-minio
run-minio:
	@ docker run -p 9000:9000 -p 9001:9001 --name minio \
       -e "MINIO_ROOT_USER=minioadmin" \
       -e "MINIO_ROOT_PASSWORD=minioadmin" \
       -v $(pwd)/data:/data \
       minio/minio server /data --console-address ":9001"
