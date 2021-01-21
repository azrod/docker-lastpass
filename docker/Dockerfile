FROM scratch

COPY --chown=65534:0 docker/docker-lastpass /
VOLUME /tmp

USER 65534

ENTRYPOINT ["/docker-lastpass"]
