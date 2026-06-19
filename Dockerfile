FROM ubuntu:latest

# 安装 ca-certificates 和 tzdata，清理 apt 缓存
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 设置架构映射和版本号
ENV SHIORI_VERSION=1.8.0

# 根据架构下载对应的 tar.gz 并解压到 /usr/bin/
RUN set -eux; \
    case "$(uname -m)" in \
        x86_64 | amd64) \
            ARCH="x86_64" \
            ;; \
        aarch64 | arm64) \
            ARCH="aarch64" \
            ;; \
        armv7l | armhf | arm) \
            ARCH="arm" \
            ;; \
        *) \
            echo "Unsupported architecture: $(uname -m)"; \
            exit 1 \
            ;; \
    esac; \
    echo "Detected architecture: $ARCH"; \
    wget -O /tmp/shiori.tar.gz "https://github.com/uparrows/shiori_cn/releases/download/${SHIORI_VERSION}/shiori_Linux_${ARCH}_${SHIORI_VERSION}.tar.gz"; \
    tar -xzf /tmp/shiori.tar.gz -C /usr/bin/; \
    rm -f /tmp/shiori.tar.gz; \
    chmod +x /usr/bin/shiori

USER root
WORKDIR /shiori
EXPOSE 8080
ENV SHIORI_DIR=/shiori/
ENTRYPOINT ["/usr/bin/shiori"]
VOLUME ["/shiori"]
CMD ["server"]
