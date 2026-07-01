# Copyright 2025-2026 Defense Unicorns
# SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

FROM alpine:3.23 AS certs
RUN apk upgrade --scripts=no apk-tools \
    && apk add --no-cache ca-certificates

FROM cgr.dev/defenseunicorns.com/busybox-fips:1.37.0

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# grab auto platform arg
ARG TARGETARCH

WORKDIR /app
RUN chown -R 65532:65532 /app && chmod -R 755 /app

# 65532 is the UID of the `nonroot` user in chainguard images.
USER 65532:65532

# copy binary from local and expose port
COPY --chown=65532:65532 build/uds-portal-linux-${TARGETARCH} /app/uds-portal
ENV PORT=8080
EXPOSE 8080

# run binary
CMD ["./uds-portal"]
