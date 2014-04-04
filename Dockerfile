FROM fedora
MAINTAINER Mrunal Patel <mpatel@redhat.com>

RUN yum install git golang -y && yum clean all

RUN export GOPATH=/root/gosrc && mkdir -p /root/gosrc/src/github.com/mrunalp && \
    cd /root/gosrc/src/github.com/mrunalp && \
    go get github.com/mrunalp/redis-todo && \
    cd /root/gosrc/src/github.com/mrunalp/redis-todo && \
    go build -o redis-todo

EXPOSE 3000

WORKDIR ["/root/gosrc/src/github.com/mrunalp/redis-todo"]

CMD ["/root/gosrc/src/github.com/mrunalp/redis-todo/redis-todo"]
