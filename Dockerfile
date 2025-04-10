FROM golang:1.23.5-alpine

RUN apk update
RUN mkdir /OCG

WORKDIR /OCG

COPY . .

RUN chmod +x /OCG/commands/start_dev.sh

CMD ["/bin/sh"]
