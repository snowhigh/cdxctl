FROM golang:1.8

LABEL io.whalebrew.name az
LABEL io.whalebrew.config.environment '["USER"]'
LABEL io.whalebrew.config.volumes '["~/.azure:/.azure"]'

RUN apt-get update && apt-get install -y libpcap-dev && apt-get install -y python-netaddr
RUN go get -v github.com/google/gopacket
RUN go get -v github.com/snowhigh/cdxctl

RUN echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" > ansible.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367
RUN apt-get update && apt-get install -y ansible

CMD ["/bin/sleep", "infinity"]
