FROM alpine

# TINI
RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]

COPY github-proxy /usr/local/bin/github-proxy

CMD ["/usr/local/bin/github-proxy"]