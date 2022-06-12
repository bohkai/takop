FROM --platform=linux/amd64 ubuntu:latest

LABEL maintainer "user"

ENV USER user
ENV HOME /home/${USER}
ENV SHELL /bin/bash

RUN apt-get update && apt-get install -y
RUN apt-get install -y git sudo curl gcc opus-tools ffmpeg

RUN curl -O https://dl.google.com/go/go1.18.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN go version

WORKDIR /home/${USER}