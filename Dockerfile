FROM alpine:3.3

RUN mkdir /gobin/

WORKDIR /gobin

EXPOSE 8080

ENTRYPOINT ["go env"]