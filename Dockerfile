FROM golang:1.20-alpine AS builder

WORKDIR /application

RUN apk add openssh gcc musl-dev make

COPY . ./

RUN make app_build

FROM alpine

WORKDIR /app

COPY --from=builder /application/app /app
COPY --from=builder /application/bin /app/bin

EXPOSE 4000

ENTRYPOINT ["bin/app"]
