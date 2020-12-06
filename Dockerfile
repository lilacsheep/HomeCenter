FROM docker.io/alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache libc6-compat
ADD build/proxy /usr/local/bin/
WORKDIR /data
EXPOSE 8080
CMD ["proxy", "-path", "/data/db"]