FROM alpine:3.3

MAINTAINER Thomas Krampl

RUN apk add --no-cache git && \
		mkdir /out

ADD generate.sh /changelog/

VOLUME /app
WORKDIR /app

CMD ["sh", "/changelog/generate.sh"]
