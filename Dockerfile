FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

ENTRYPOINT ["/app/bin/mebot"]