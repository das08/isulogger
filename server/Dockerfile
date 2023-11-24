FROM golang:1.18.0-alpine

WORKDIR /go/src/app
#COPY go.mod go.sum ./
#RUN go mod download && go mod verify

# install gcc, which is required by go-sqlite3
ENV CGO_ENABLED=1
ENV DEBIAN_FRONTEND=noninteractive
RUN apk add --no-cache git sqlite-dev gcc musl-dev

CMD ["go", "run", "."]