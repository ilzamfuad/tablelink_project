# tablelink_project

This is for Tablelink Test Code Project, The test include grpc connection with restful API, and User Service API

## Requirements

- **Golang**: Latest version
- **Protobuf**: Latest version

## Installation

1. Install dependencies:
  ```bash
  git clone https://github.com/googleapis/googleapis.git
  ```
  Ensure the dependencies are placed in the `dependencies` folder.

## Development

To add or update APIs in the `.proto` files, use the following command:
```bash
protoc -I=. -I=dependencies/googleapis --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/<proto_file>
```
Replace `<proto_file>` with the name of your `.proto` file.

## Running the Application

To run the application, use the following command:
```bash
go run server/main.go
```
