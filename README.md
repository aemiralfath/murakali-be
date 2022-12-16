# Murakali - BE

## How to run
1. install go 1.19
2. install docker desktop & docker compose
3. install golang migrate (https://github.com/golang-migrate/migrate. optional, only need for doing migrations)
4. create .env file and copy the value from confluence page (https://murakali.atlassian.net/wiki/spaces/M/pages/1474562/.env)
5. run `go mod tidy` in terminal
6. run `make docker-up` to start the server
7. you can access the BE server on `http://localhost:8080/`
8. open `http://localhost:8081/` and login using credentials to run sql seeder command
9. read makefile command to understand other command
