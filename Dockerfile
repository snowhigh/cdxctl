FROM golang:1.8

LABEL io.whalebrew.name az
LABEL io.whalebrew.config.environment '["USER"]'
LABEL io.whalebrew.config.volumes '["~/.azure:/.azure"]'

RUN apt-get update && apt-get install -y libpcap-dev python-netaddr sshpass python-pip python-dev build-essential libssl-dev libffi-dev jq
RUN go get -v github.com/google/gopacket
RUN go get -v github.com/snowhigh/cdxctl

RUN pip install --upgrade cffi
RUN pip install ansible

CMD ["/bin/sleep", "infinity"]
