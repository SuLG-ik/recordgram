FROM golang:1.19.3-alpine as builder

WORKDIR /app/go-sample-app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/app .

FROM alpine:3.17

WORKDIR /run

COPY --from=builder /app/go-sample-app/out/app /recordgram/

CMD ["/recordgram/app"]