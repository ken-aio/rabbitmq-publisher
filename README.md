# rabbitmq-publisher
Bulk publish messages to rabbitmq.  
```
$ go run main.go --uri=amqp://localhost:5672/vhost1 --routing-key=hello.world --exchange=hello --file-path=./test-file.tsv
2021/01/19 14:12:48 success to publish:  {"message": "hello world"}
2021/01/19 14:12:48 success to publish:  {"message": "hello world part2"}
2021/01/19 14:12:48 success to publish:  {"message": "hello world part3"}
```

```
$ cat test-file.tsv
{"message": "hello world"}
{"message": "hello world part2"}
{"message": "hello world part3"}
```
