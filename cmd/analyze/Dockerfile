FROM golang:1.21.6@sha256:7b575fe0d9c2e01553b04d9de8ffea6d35ca3ab3380d2a8db2acc8f0f1519a53 as build
RUN apt-get update && apt-get install -y libpcap-dev
WORKDIR /src

# First cache the dependencies to avoid downloading again on code change
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . ./
RUN go build -o analyze cmd/analyze/main.go && go build -o worker cmd/worker/main.go

FROM ubuntu:22.04@sha256:e6173d4dc55e76b87c4af8db8821b1feae4146dd47341e4d431118c7dd060a74

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        iptables \
        iproute2 \
        podman \
        software-properties-common && \
    update-alternatives --set iptables /usr/sbin/iptables-legacy && \
    update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy

# Install gVisor.
RUN mkdir -m 0700 -p /etc/apt/keyrings && \
    curl -fsSL https://gvisor.dev/archive.key -o /etc/apt/keyrings/gvisor.key && \
    echo "deb [signed-by=/etc/apt/keyrings/gvisor.key] https://storage.googleapis.com/gvisor/releases 20220425 main" > /etc/apt/sources.list.d/gvisor.list && \
    apt-get update && apt-get install -y runsc

COPY --from=build /src/analyze /usr/local/bin/analyze
COPY --from=build /src/worker /usr/local/bin/worker
COPY --from=build /src/tools/gvisor/runsc_compat.sh /usr/local/bin/runsc_compat.sh
COPY --from=build /src/tools/network/iptables.rules /usr/local/etc/iptables.rules
COPY --from=build /src/tools/network/podman-analysis.conflist /etc/cni/net.d/podman-analysis.conflist
RUN chmod 755 /usr/local/bin/runsc_compat.sh && \
    chmod 644 /usr/local/etc/iptables.rules /etc/cni/net.d/podman-analysis.conflist

ARG SANDBOX_IMAGE_TAG
ENV OSSF_SANDBOX_IMAGE_TAG=${SANDBOX_IMAGE_TAG}
