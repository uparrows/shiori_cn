# build stage
FROM docker.io/golang:1.21.7-alpine3.19 AS builder
WORKDIR /src
COPY . .
RUN go build -ldflags '-s -w'

# server image

FROM docker.io/alpine:3.19
LABEL org.opencontainers.image.source https://github.com/go-shiori/shiori
COPY --from=builder /src/shiori /usr/bin/
USER root
WORKDIR /shiori
EXPOSE 8080
ENV SHIORI_DIR /shiori/
ENTRYPOINT ["/usr/bin/shiori"]
CMD ["server"]
