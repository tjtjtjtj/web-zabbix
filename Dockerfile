From golang:1.9.2-alpine3.7

RUN apk --update add \
    chromium \
    chromium-chromedriver \
    udev \
    git

WORKDIR /noto

RUN wget https://noto-website.storage.googleapis.com/pkgs/NotoSansCJKjp-hinted.zip && \
    unzip NotoSansCJKjp-hinted.zip && \
    mkdir -p /usr/share/fonts/noto && \
    cp *.otf /usr/share/fonts/noto && \
    chmod 644 -R /usr/share/fonts/noto/ && \
    fc-cache -fv && \
    rm -rf /noto

WORKDIR /go/src/github.com/tjtjtjtj/web-zabbix

RUN go get github.com/sclevine/agouti

#CMD go run *.go
