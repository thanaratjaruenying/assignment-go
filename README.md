# How to run the server
Clone the repository and run the following command:

First, you need to install the dependencies.

```bash
$ go mod download
```
then create a `.env` file in the root directory of the project as below.
    
```
PORT=3000
```
after that, you can run the server.
    
```bash
$ go run main.go
```

If you want to run the server in container, you can use the following command:
    
```bash
$ docker build -t go-server .
```
then
```bash
$ docker run -it -p PORT:PORT go-server
```