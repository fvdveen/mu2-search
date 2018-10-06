FROM golang:1.11 AS build

WORKDIR /mu2-search

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN CGO_ENABLED=0 go build -o mu2-search main.go

FROM alpine:latest AS RUN

WORKDIR /app/
COPY --from=build /mu2-search/mu2-search mu2-search
RUN apk add --no-cache ca-certificates

CMD ["./mu2-search"]
