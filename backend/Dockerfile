FROM golang:alpine

RUN apk update && \
    apk add --no-cache ffmpeg

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8000

CMD ["go", "run", "main.go" ]