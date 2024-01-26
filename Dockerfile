FROM golang:alpine as build

RUN apk update && apk upgrade && apk add --no-cache automake make gettext
WORKDIR /app
COPY . ./
#     apk add --no-cache bash git openssh
RUN make build

CMD ["/app/run.sh"]

