FROM gcr.io/distroless/static:nonroot
COPY swag2mcp /usr/local/bin/swag2mcp
USER 65532:65532
ENTRYPOINT ["swag2mcp"]
