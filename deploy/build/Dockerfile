FROM golang:1.15-alpine as stage-build
LABEL stage=stage-build
WORKDIR /build/kubeops
ARG GOARCH

ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=$GOARCH
ENV CGO_ENABLED=0


RUN  apk update \
  && apk add git \
  && apk add make
COPY ../../go.mod go.sum ./
RUN go mod download
COPY ../.. .
RUN make build_server_linux GOARCH=$GOARCH

FROM alpinelinux/ansible:latest

RUN apk add sshpass \
    && apk add rsync \
    && apk add openssl \
    && pip3 install netaddr \
    && pip3 install pywinrm


RUN mkdir /root/.ssh  \
    && touch /root/.ssh/config \
    && echo -e "Host *\n\tStrictHostKeyChecking no\n\tUserKnownHostsFile /dev/null" > /root/.ssh/config

COPY --from=stage-build /build/kubeops/dist/etc /etc/
COPY --from=stage-build /build/kubeops/dist/usr /usr/
COPY --from=stage-build /build/kubeops/dist/var /var/

RUN echo 'kubeops-server' >> /root/entrypoint.sh

VOLUME ["/var/kubeops/data"]

EXPOSE 8080

CMD ["sh","/root/entrypoint.sh"]
