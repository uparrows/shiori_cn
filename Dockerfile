# build stage
FROM golang:1.23.6-alpine3.21 AS builder
WORKDIR /src
COPY . .
RUN go build -ldflags '-s -w'

# server image

FROM alpine:3.21
LABEL org.opencontainers.image.source https://github.com/go-shiori/shiori
COPY --from=builder /src/shiori /usr/bin/
RUN apk add --no-cache ca-certificates tzdata
USER root
WORKDIR /shiori
EXPOSE 8080
ENV SHIORI_DIR /shiori/
ENTRYPOINT ["/usr/bin/shiori"]
VOLUME [ "/shiori" ]
CMD ["server"]