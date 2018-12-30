EXECPATH=build/exec
vendor: 
	go mod vendor
build: 
	go build -mod=vendor -o=$(EXECPATH)
exec:
	./$(EXECPATH)
run: vendor build exec