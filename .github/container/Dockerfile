FROM alpine

# TINI
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

COPY github-proxy /usr/local/bin/github-proxy

CMD ["/usr/local/bin/github-proxy"]