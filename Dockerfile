FROM docker.io/alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache libc6-compa
ADD build/proxy /usr/local/bin/
WORKDIR /data

CMD ["proxy", "-path", "/data/db"]