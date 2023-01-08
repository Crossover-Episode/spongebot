FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /spongebot ./main.go

FROM golang:1.18-alpine

RUN apk --no-cache add ca-certificates

COPY --from=build /spongebot /spongebot
ENTRYPOINT ["/spongebot"]