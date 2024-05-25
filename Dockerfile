FROM alpine:3.20.0 AS base

RUN apk update
RUN apk upgrade
RUN apk add --update go=1.22.3-r0 

FROM base AS tester

WORKDIR /opt/url-short

ADD . /opt/url-short

CMD ["go", "test"]

FROM base AS builder

WORKDIR /build

ADD . /build

RUN go build -o main .

FROM builder AS production

WORKDIR /opt/url-short/

COPY --from=builder /build/main .

CMD ["./main"]
