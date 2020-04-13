FROM golang:1.14.2-buster AS builder

WORKDIR /frpc

# Install dependencies first, so we can cache them.
RUN go get github.com/go-chi/chi
RUN go get github.com/mitchellh/go-homedir
RUN go get github.com/mitchellh/go-wordwrap
RUN go get github.com/stretchr/testify

COPY . .

RUN go install ./...

FROM factoriotools/factorio:0.17.79

COPY --from=builder /go/bin/frpc-sidecar /usr/local/bin/frpc-sidecar
COPY --from=builder /frpc/mod /factorio/mods/frpc_0.0.1

ENTRYPOINT /bin/bash -c "frpc-sidecar -dir /factorio/script-output > frpc-sidecar.stdout 2> frpc-sidecar.stderr & exec /docker-entrypoint.sh"
