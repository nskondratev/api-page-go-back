# API Page backend
This is the backend part of API Page application (documentation + real-time Socket.io API testing).

## Prerequisites
* Go v1.12.4
* MySQL v5.7.23

## Build
```bash
go build
```

## Run
### Development
Start application locally:
```bash
go run main.go
```

### Production
* Build the app and upload binaries to server
* Upload example service file to server and modify it
* Run application as service:
```bash
sudo systemctl start api-page-backend
```

## API Docs
* [Websockets API](docs/WS_API.md)
