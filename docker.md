## Docker commands
```sh
#Up containers
make compose-db-up
#Up migration
make up
#Recreate mocks
mockery
#Run integration test
make test
#Run program
make run
#Remove container
make compose-db-rm
#cd to test file dir
go test ./... -coverprofile=coverage.out
go tool cover -html coverage.out -o coverage.html
#Open coverage.html using windows File Explorer
```


