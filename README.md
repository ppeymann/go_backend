## Example backend API Service
#### -This is a backend project for practicing frontend developers who need a real backend

### - Clone Repo:
```git clone https://github.com/ppeymann/go_backend.git```

### - Build docker Image:
```docker build --tag example:latest .```

### - Run:
```cd cmd/eg```

```go run main.go```

### - Generating swagger docs:
```swag init --parseDependency --parseInternal -g /server/server.go```


