FROM scratch
COPY --chown=65534:65534 docker/docker-lastpass /
USER 65534
VOLUME /tmp

ENTRYPOINT ["/docker-lastpass"]