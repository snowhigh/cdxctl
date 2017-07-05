FROM golang:1.8

RUN apt-get update && apt-get install -y libpcap-dev python-netaddr sshpass python-pip python-dev build-essential libssl-dev libffi-dev jq vim nginx net-tools
RUN go get -v github.com/google/gopacket
RUN go get -v github.com/simonschuang/cdxctl

RUN pip install --upgrade cffi
RUN pip install ansible ansible-cmdb

RUN wget https://storage.googleapis.com/kubernetes-release/release/v1.5.7/bin/linux/amd64/kubectl -O /usr/local/sbin/kubectl && chmod +x /usr/local/sbin/kubectl

RUN mkdir -p /etc/ansible
ADD ansible.cfg /etc/ansible/ansible.cfg

CMD ["/bin/sleep", "infinity"]
