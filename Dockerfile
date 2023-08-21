FROM golang:alpine

# Install OS-level dependencies.
RUN apk add --no-cache curl git
RUN apk add --no-cache ca-certificates

# Copy our source code into the container.
WORKDIR /go/src/github.com/chihaya/chihaya
ADD ./chihaya /go/src/github.com/chihaya/chihaya

# Install our golang dependencies and compile our binary.
RUN CGO_ENABLED=0 go install ./cmd/chihaya

RUN adduser -D chihaya

# Expose a docker interface to our binary.
EXPOSE 8080

# Drop root privileges
USER chihaya

ADD config.yaml /etc/chihaya.yaml

CMD ["/go/bin/chihaya"]