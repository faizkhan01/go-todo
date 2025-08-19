FROM golang:1.20-alpine AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o todo-app

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/todo-app /app/todo-app
EXPOSE 8080
CMD ["/app/todo-app"]