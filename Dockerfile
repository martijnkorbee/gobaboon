FROM golang:1.20-alpine AS builder

WORKDIR /app

RUN apk add openssh gcc musl-dev make

COPY . ./

RUN make app_build

FROM alpine

WORKDIR /app
COPY --from=builder /app /app

EXPOSE 4000

CMD ["./bin/app"]

