FROM alpine:latest
RUN apk add --no-cache su-exec
ENV HOME=/home/nonroot
COPY swag2mcp /usr/local/bin/swag2mcp
COPY scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
