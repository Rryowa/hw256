## Docker commands
```sh
#Up db containers
make compose-db-up
#Up migration & build
make all
#Recreate mocks (optional)
mockery
#Run integration test
make test
#Run program
make run
#Down migrations
make down
#Remove container
make compose-db-rm
#cd to test file dir
go test ./... -coverprofile=coverage.out
go tool cover -html coverage.out -o coverage.html
#Open coverage.html using windows File Explorer
```


