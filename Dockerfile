FROM golang:1.25.0-alpine3.22 AS builder

ARG APP_VERSION="dev"
ARG APP_BUILD="dev"

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build

FROM alpine:3.22.1

WORKDIR /app
COPY --from=builder /app/ai-instructions /usr/local/bin/
COPY --from=builder /app/rules ./rules

CMD ["help"]
ENTRYPOINT ["ai-instructions"]