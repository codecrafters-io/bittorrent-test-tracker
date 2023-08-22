FROM golang:alpine

# Install OS-level dependencies.
RUN apk add --no-cache curl git
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache transmission-cli

# Copy our source code into the container.
WORKDIR /go/src/github.com/chihaya/chihaya
ADD ./chihaya /go/src/github.com/chihaya/chihaya

# Install our golang dependencies and compile our binary.
RUN CGO_ENABLED=0 go install ./cmd/chihaya

# Expose a docker interface to our binary.
EXPOSE 8080

ADD config.yaml /etc/chihaya.yaml

ADD tracker.sh /bin/tracker.sh
ADD seeder.sh /bin/seeder.sh
ADD torrent_files /etc/torrent_files

CMD ["/bin/tracker.sh"]