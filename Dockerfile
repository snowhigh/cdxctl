FROM golang:1.8

RUN apt-get update && apt-get install -y libpcap-dev python-netaddr sshpass python-pip python-dev build-essential libssl-dev libffi-dev jq
RUN go get -v github.com/google/gopacket
RUN go get -v github.com/simonschuang/cdxctl

RUN pip install --upgrade cffi
RUN pip install ansible ansible-cmdb

CMD ["/bin/sleep", "infinity"]
