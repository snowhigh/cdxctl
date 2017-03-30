FROM golang:1.8

LABEL io.whalebrew.name az
LABEL io.whalebrew.config.environment '["USER"]'
LABEL io.whalebrew.config.volumes '["~/.azure:/.azure"]'

RUN apt-get update && apt-get install -y libpcap-dev && apt-get install -y python-netaddr
RUN go get -v github.com/google/gopacket
RUN go get -v github.com/snowhigh/cdxctl

CMD ["/bin/sleep", "infinity"]
