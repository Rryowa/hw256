## Docker commands
```sh
#Up containers
make compose-up
#Up migration
make up
#Recreate mocks
mockery
#Run integration test
make test
#Run program
make run
#Remove container
make compose-rm
#cd to test file dir
go test -coverprofile=coverage.out
go tool cover -func=coverage.out
```


