# build binary first
FROM    figassis/ubuntu-golang
LABEL   maintainer="Assis Ngolo <figassis@gmail.com>"

ENV LOG_LEVEL=info
ENV CONFIG_FILE=/etc/mysql-backup/config.yaml

# add glide config and install dependencies with glide in a separate step to speed up subsequent builds
WORKDIR /go/src/github.com/figassis/mysql-backup

# add source and build package
ADD . /go/src/github.com/figassis/mysql-backup/
RUN apt-get update && apt-get -y install mysql-client bzip2 curl upx \
    && wget https://github.com/restic/restic/releases/download/v0.9.6/restic_0.9.6_linux_amd64.bz2 \
    && bzip2 -d restic_0.9.6_linux_amd64.bz2 && mv restic_0.9.6_linux_amd64 /usr/local/bin/restic \
    && chmod +x /usr/local/bin/restic && restic self-update \
    && mkdir -p /etc/mysql-backup

RUN go mod tidy \
    # https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
    && go build -i -o /mysql-backup -ldflags="-s -w" \
    # compress binary
    && upx /mysql-backup

ENTRYPOINT ["/mysql-backup"]
CMD        [""]
