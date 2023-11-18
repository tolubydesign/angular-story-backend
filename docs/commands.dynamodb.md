
[AWS dynamodb CLI Reference](https://docs.aws.amazon.com/cli/latest/reference/dynamodb/)

__Delete Table__
```sh
$ aws dynamodb delete-table --table-name Story --endpoint-url http://localhost:8000
```

__List Tables in Database__
```sh
$ aws dynamodb list-tables --endpoint-url http://localhost:8000
```

__Scan Items in Database__
```sh
$ aws dynamodb scan --endpoint-url http://localhost:8000 --table-name Story
```