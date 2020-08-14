build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o ./consul-test -v main.go
	docker build -t viniciusramosdefaria/consul-test:latest .
	docker push viniciusramosdefaria/consul-test:latest