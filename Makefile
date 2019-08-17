.PHONY: mod clean build start

mod:
	go get -u ./...

clean: 
	rm -rf ./nature-remo/nature-remo
	
build:
	GOOS=linux GOARCH=amd64 go build -o nature-remo/nature-remo ./nature-remo

start: build ## Start Lambda on localhost
	sam local start-api --env-vars env.json

create-table:
	aws dynamodb create-table --cli-input-json file://ddl/create_remo_table.json --endpoint-url http://localhost:18000

insert-data:
	aws dynamodb batch-write-item --request-items file://ddl/put_remo_sample_date.json --endpoint-url http://localhost:18000

delete-table:
	aws dynamodb delete-table --endpoint-url http://localhost:18000 --table-name Remo