# Todo App (Go + Gin)

Simple in-memory TODO REST API written in Go using the Gin framework.

## Requirements

- Go 1.18+
- Git (optional)

## Run locally

```bash
go mod download
go run main.go
```

Server runs at `http://localhost:8080`.

## API

- `GET /todos` - list todos
- `GET /todos/:id` - get todo by id
- `POST /todos` - create todo `{ "title": "Task", "details": "..." }`
- `PUT /todos/:id` - update todo (partial) `{ "title": "...", "done": true }`
- `DELETE /todos/:id` - delete todo

## Docker

Build:
```
docker build -t todo-app .
docker run -p 8080:8080 todo-app
```

## Notes

- Data is stored in memory (non-persistent).