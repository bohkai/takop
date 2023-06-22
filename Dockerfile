FROM --platform=linux/arm64 ubuntu:latest

LABEL maintainer "user"

ENV USER user
ENV HOME /home/${USER}
ENV SHELL /bin/bash
WORKDIR /home/${USER}

ENV MECAB_VERSION 0.996
ENV IPADIC_VERSION 2.7.0-20070801
ENV MECAB_URL https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7cENtOXlicTFaRUE
ENV IPADIC_URL https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7MWVlSDBCSXZMTXM
ENV CGO_LDFLAGS="-L/usr/local/Cellar/mecab/0.996/lib -lmecab -lstdc++ -lm"
ENV CGO_CFLAGS="-I/usr/local/Cellar/mecab/0.996/include"

RUN apt-get update && apt-get install -y
RUN apt-get install -y git sudo curl gcc opus-tools ffmpeg make g++ xz-utils file pkg-config libopus-dev libopusfile-dev
RUN curl -O https://dl.google.com/go/go1.20.5.linux-arm64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.5.linux-arm64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN go version

RUN curl -SL -o mecab-${MECAB_VERSION}.tar.gz ${MECAB_URL} && \
    tar zxf mecab-${MECAB_VERSION}.tar.gz && \
    cd mecab-${MECAB_VERSION} && \
    ./configure --enable-utf8-only --with-charset=utf8 --build=arm && \
    make && \
    make install

RUN cd && \
    curl -SL -o mecab-ipadic-${IPADIC_VERSION}.tar.gz ${IPADIC_URL} && \
    tar zxf mecab-ipadic-${IPADIC_VERSION}.tar.gz && \
    cd mecab-ipadic-${IPADIC_VERSION} && \
    ./configure --with-charset=utf8 --build=arm && \
    make && \
    make install

RUN cd && \
    git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git /mecab-ipadic-neologd && \
    /mecab-ipadic-neologd/bin/install-mecab-ipadic-neologd -n -a -y -p /usr/local/lib/mecab/dic/mecab-ipadic-neologd && \
    echo "dicdir = /usr/local/lib/mecab/dic/mecab-ipadic-neologd" > /usr/local/etc/mecabrc

RUN rm -rf \
    mecab-${MECAB_VERSION}* \
    mecab-${IPADIC_VERSION}* \
    /mecab-ipadic-neologd

WORKDIR /home/${USER}/takop
COPY . .

CMD [ "make", "run" ]
