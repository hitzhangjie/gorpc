FROM hitzhangjie/centos:latest

MAINTAINER hitzhangjie

ENV GIT_SSL_NO_VERIFY=true GIT_TERMINAL_PROMPT=1
ENV PATH=$PATH:/usr/local/go/bin

USER root

COPY bash/bashrc /root/.bashrc
COPY git/gitconfig /etc/gitconfig
COPY vim.tgz /root/vim.tgz
COPY protobuf /root/protobuf

RUN echo "===> config vim"                                                  && \
    tar xvfz /root/vim.tgz -C /root                                         && \
    rm /root/vim.tgz                                                        && \
    echo "===> config yum"                                                  && \
    sed -i 's/tsflags=nodocs/#tsflags=nodocs/' /etc/yum.conf                && \
    rm -rf /etc/yum.repos.d && mkdir /etc/yum.repos.d                       && \
    curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.cloud.tencent.com/repo/centos7_base.repo && \
    echo "===> install common utilities"                                    && \
    yum install -y sudo man-db man-pages make which git vim tree gcc ctags-etags radare2 && \
    yum install -y /root/protobuf/*.rpm                                     && \
    rm -rf /root/protobuf                                                   && \
    yum clean all                                                           && \
    echo "===> config golang"                                               && \
    curl -o /root/go.tgz https://dl.google.com/go/go1.16.13.linux-amd64.tar.gz && \
    tar xvfz /root/go.tgz -C /usr/local                                     && \
    rm /root/go.tgz                                                         && \
    echo "===> config golang utilities"                                     && \
    go get github.com/go-delve/delve/cmd/dlv                                && \
    go get github.com/golang/protobuf/protoc-gen-go                         && \
    go clean -cache

WORKDIR /root
