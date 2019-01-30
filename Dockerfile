FROM golang:1.11-alpine  AS build_base

# Force the go compiler to use modules
RUN apk add --no-cache git openssh-client musl-dev gcc
ENV GO111MODULE=on
WORKDIR github.com/alekc/docker-dbg

COPY go.mod .
COPY go.sum .

#Download all dependencies.
RUN go mod download

FROM build_base AS server_builder

#copy the source
COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a .

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine AS final-build
RUN apk add --no-cache tzdata ca-certificates

EXPOSE 80

# Finally we copy the statically compiled Go binary.
COPY --from=server_builder /go/bin/docker-dbg /bin/docker-dbg

ENTRYPOINT ["/bin/docker-dbg"]
