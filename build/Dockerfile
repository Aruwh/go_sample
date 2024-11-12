####################
# BUILD IMAGE      #
####################

FROM alpine:latest AS builder
RUN echo `pwd`

RUN apk update && apk add --no-cache ca-certificates

RUN mkdir /user && \
  echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
  echo 'nobody:x:65534:' > /user/group

WORKDIR /app


COPY ./bin/fewoserv .

####################
# PRODUCTION IMAGE #
####################
FROM alpine:latest AS final

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/fewoserv .

USER nobody:nobody

ENTRYPOINT ["./fewoserv"]
