# iacgen
iacgen generates Iac (Infrastructure as Code) for the given input configuration

## Run locally
```sh
$ make build
$ ./bin/iac-gen
```
## Run docker locally
```sh
$ docker build -t iacgen -f Dockerfile .
$ docker run -p 8000:8000 --rm --name iac-gen iacgen:latest
```
## Example
```sh
$ curl -i --output terraform.zip -X POST localhost:8000/v1/api/iac/generate -d '{"metadata":{"repository":"xyz.git","commit":"xyz","branch":"xyz"},"graph":{"test1":["test2","test3"],"test2":["test3"],"test3":["test1"]},"projects":[{"metadata":{"location":"./test1","name":"test1","languages":["java"],"frameworks":["spring"]},"services":[{"name":"test1","dns":"test1","image":"test/image","ports":[{"listen":8080,"forward":80},{"listen":8081}],"http_egress":[{"endpoint":"http://test2/foo","operations":["GET","POST"]},{"endpoint":"http://test5/foo","operations":["GET"]}],"http_ingress":[{"endpoint":"/foo","operations":["POST"]}],"s3":[{"bucket":"xyz","operations":["GET","PUT","DELETE"]}],"database":[{"dsn":"jdbc:mysql://mysql.db.server:3306/my_database"}]},{"name":"test2","dns":"test2","ports":[{"listen":8080,"forward":80}],"http_egress":[{"endpoint":"http://test1/foo","operations":["GET"]}],"s3":[{"bucket":"xyz"}],"database":[{"dsn":"jdbc:mysql://mysql.db.server:3306/my_database_2"}]}]}]}'
$ unzip terraform.zip
```
