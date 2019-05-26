# Domainfo / Golang
## It's running right now!
This API is running [here](https://domainfo-mariopayan.c9users.io/v1/api/domain/list) and it's using 
* AWS Cloud9 as server
* AWS EC2 cockroach instance as DB 

Or continue in order to make it running locally

## Dependencies
* go1.12.5
* cockroach-v19.1.1

## Install
* Go to project directory
* Run the command `go get -d ./...` to install all dependencies

## Run

### Database
* Open a new console window
* Run the command `cockroach start --insecure --listen-addr=localhost`
* Leave that console window there and go back to the main console window
* Run the commmand `cockroach sql --insecure --host=localhost --port=26257 --database=domainfo < init.sql` inside the project folder

### API
* Go to project directory
* Run the command `go run *.go`

## Usage
You will see the endpoints availables in the console, the default url is `localhost:3333`