# build binary first
FROM    figassis/ubuntu-golang
LABEL   maintainer="Assis Ngolo <figassis@gmail.com>"

ENV LOG_LEVEL=info
ENV CONFIG_FILE=/etc/mysql-backup/config.yaml

# add glide config and install dependencies with glide in a separate step to speed up subsequent builds
WORKDIR /go/src/github.com/figassis/mysql-backup

# add source and build package
ADD . /go/src/github.com/figassis/mysql-backup/
RUN apt-get update && apt-get -y install restic mysql-client bzip2 curl \
    && mkdir -p /etc/mysql-backup

RUN go mod tidy \
    # https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
    && go build -i -o /mysql-backup -ldflags="-s -w" \
    # compress binary
    && upx /mysql-backup

ENTRYPOINT ["/mysql-backup"]
CMD        [""]
