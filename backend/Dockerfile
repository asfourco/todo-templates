FROM golang:latest as backend
EXPOSE 8080
ARG BUILD_FLAGS
ENV BUILD_FLAGS $BUILD_FLAGS
ENV DEBUG=.+

ENV TINI_VERSION "v0.19.0"


COPY . /opt/app
COPY go.mod /opt/app
COPY go.sum /opt/app
COPY Makefile /opt/app

WORKDIR /opt/app

RUN set -eux && make install


# Add Tini
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

CMD ["/tini", "--", "todo", "serve"]
