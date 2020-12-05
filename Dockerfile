FROM docker.io/alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
ADD build/proxy /usr/local/bin/
WORKDIR /data

CMD ["proxy", "-path", "/data/db"]