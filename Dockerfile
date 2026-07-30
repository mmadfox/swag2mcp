FROM alpine:latest
RUN apk add --no-cache su-exec && \
    adduser -D -u 65532 -h /home/nonroot nonroot
ENV HOME=/home/nonroot
COPY swag2mcp /usr/local/bin/swag2mcp
COPY scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
