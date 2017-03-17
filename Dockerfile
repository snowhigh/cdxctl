FROM golang:1.8

LABEL io.whalebrew.name az
LABEL io.whalebrew.config.environment '["USER"]'
LABEL io.whalebrew.config.volumes '["~/.azure:/.azure"]'

RUN go get github.com/google/gopacket
RUN apt-get update && apt-get install -y libpcap-dev

ENTRYPOINT ["echo hello!"]
