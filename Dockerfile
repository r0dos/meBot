FROM golang:alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
    		-o ./bin/mebot \
    		./cmd/mebot

FROM alpine

WORKDIR /app

COPY --from=builder /src/bin/mebot .
COPY --from=builder /src/config.yml .
#COPY --from=builder /src/locales locales
#COPY --from=builder /src/sql sql

ENTRYPOINT ["/app/mebot"]