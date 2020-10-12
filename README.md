1. db connection is in store package(3rd party cloud postgres instance)
2. rabbitMQ connection is in rabbitMQ package(3rd party cloud instance)
3. change connections according to requirment then run 'main.go' to start the server
4. unit test file added to service package (run command 'go test -v')
5. publishing is done via post restful api (endpoint : 'localhost/8080/offer')
6. consumer will automatically start once it recieve a task on queue and save filtered data to database.
