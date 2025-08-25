FROM golang:1.24.6 AS builder

WORKDIR /src

COPY go.mod go.sum ./

# Cache dependencies
RUN go mod download
RUN go mod verify

# Once we have `--parents` in the next docker version
# this can just be:
# COPY --parents ./**/*.go ./

COPY "./cmd/" "./cmd/"
COPY internal/ ./internal/
COPY main.go ./

ENV CGO_ENABLED=0

RUN go build .

FROM scratch AS runner

COPY --from=builder /src/micro-nas-lcd /

ENTRYPOINT ["/micro-nas-lcd"]
CMD ["display"]