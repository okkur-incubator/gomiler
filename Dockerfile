FROM alpine:3.9

RUN apk --no-cache add ca-certificates

ADD gomiler /gomiler

CMD ["/gomiler"]