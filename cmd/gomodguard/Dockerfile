FROM gcr.io/distroless/static
ARG TARGETOS
ARG TARGETARCH
COPY ${TARGETOS}/${TARGETARCH}/gomodguard /gomodguard
ENTRYPOINT ["/gomodguard"]
