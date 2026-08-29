# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------------
# Stage 1: build
#
# --platform=$BUILDPLATFORM keeps the compiler running natively while
# cross-compiling to $TARGETPLATFORM, which is free here because the build is
# cgo-free: the SQLite driver is modernc.org/sqlite (pure Go), and PostgreSQL
# goes through pgx.
# ---------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Dependency layer: only invalidated when go.mod/go.sum change.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

COPY . .

# Static assets (public/static/main.css) are embedded into the binary by
# public/fs.go, so run `make css` before building this image if the Tailwind
# sources changed — nothing here regenerates them.
#
# -tags timetzdata embeds the zoneinfo database, which scratch does not carry.
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -trimpath -tags timetzdata -ldflags="-s -w" -o /out/web ./cmd/web

# scratch has no shell, so anything the app expects to exist at runtime has to
# be staged here and copied in:
#   /tmp          net/http spills multipart uploads over 32 MiB to os.TempDir()
#   /app/dbs      SQLite task queue (openSQLite creates the file, not the mount)
#   /app/uploads  afero base path for stored files
RUN mkdir -p /out/rootfs/tmp /out/rootfs/app/dbs /out/rootfs/app/uploads \
  && chmod 1777 /out/rootfs/tmp \
  && echo 'app:x:10001:10001::/app:/sbin/nologin' > /out/passwd \
  && echo 'app:x:10001:' > /out/group

# ---------------------------------------------------------------------------
# Stage 2: runtime
# ---------------------------------------------------------------------------
FROM scratch

# Outbound HTTPS (object storage, SMTP over TLS) needs the trust store.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/passwd /etc/passwd
COPY --from=builder /out/group /etc/group
COPY --from=builder --chown=10001:10001 /out/rootfs/ /
COPY --from=builder /out/web /usr/local/bin/

# viper resolves ./config/config.yaml relative to the working directory and
# fails hard if it is missing. Every value is overridable at runtime with
# GO_-prefixed env vars (GO_DATABASE_CONNECTION, GO_APP_HOST,
# GO_APP_ENCRYPTIONKEY, GO_APP_ENVIRONMENT, ...) — do not bake
# environment-specific values into the image.
COPY --chown=10001:10001 config/config.yaml /app/config/config.yaml

WORKDIR /app
USER 10001:10001

# Serves plain HTTP only (no TLS terminated in-process); matches http.port in
# config/config.yaml. Terminate TLS at a proxy/load balancer in front.
EXPOSE 8000

# The task queue and stored files are state; without volumes they are lost on
# every redeploy.
VOLUME ["/app/dbs", "/app/uploads"]

# No HEALTHCHECK: scratch has no shell or wget to exec. Probe from outside
# (compose `depends_on`, or an orchestrator HTTP/TCP probe against :8000).
ENTRYPOINT ["/usr/local/bin/web"]
