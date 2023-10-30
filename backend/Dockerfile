FROM golang:alpine

RUN apk update && \
    apk add --no-cache ffmpeg

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build main.go

EXPOSE 8000

CMD [ "./main" ]