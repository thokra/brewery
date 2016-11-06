FROM alpine:3.3

MAINTAINER Thomas Krampl

RUN apk add --no-cache git && \
		mkdir /out

ADD init.sh /usr/local/bin/

CMD ["/bin/sh", "-c", "/usr/local/bin/init.sh"]
