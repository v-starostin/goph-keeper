.PHONY: run
run:
	export DB_CONN_STRING=.db_conn && export SECRET=.secret && go run internal/config/*.go
