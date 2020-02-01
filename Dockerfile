FROM alpine:3.8

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ADD version.txt /etc
ADD ./bin/go-utils /bin/

RUN chmod +x /bin/go-utils

CMD ["/bin/go-utils"]