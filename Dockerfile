# build stage
FROM ghcr.io/ghcri/golang:1.17-alpine3.15 AS builder
WORKDIR /src
COPY . .
RUN go build -v -ldflags '-s -w'

# server image

FROM docker.io/alpine:3.15
LABEL org.opencontainers.image.source https://github.com/uparrows/shiori_cn
COPY --from=builder /src/shiori /usr/bin/
USER root
WORKDIR /shiori
EXPOSE 8080
ENV SHIORI_DIR /shiori/
ENTRYPOINT ["/usr/bin/shiori"]
CMD ["server"]
