FROM debian:12

ENV container=docker

RUN apt update && apt install -y \
    wget \
    curl \
    openssl \
    systemd \
    systemd-sysv \
    sudo \
    ufw \
    iproute2 \
    procps \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && ARCH=$(dpkg --print-architecture) \
    && wget -q "https://go.dev/dl/go1.23.6.linux-${ARCH}.tar.gz" -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/root/go"
ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR /app

CMD ["/sbin/init"]