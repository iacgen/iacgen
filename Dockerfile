# -------- builder stage -------- #
FROM golang:1.19 AS builder

ARG GOOS_VAL=linux
ARG GOARCH_VAL=amd64
ARG CGO_ENABLED_VAL=0

WORKDIR $GOPATH/src/iac-gen
COPY . ./


# build binary
RUN CGO_ENABLED=${CGO_ENABLED_VAL} GOOS=${GOOS_VAL} GOARCH=${GOARCH_VAL} \
    go build -v -ldflags "-w -s -X main.VERSION=${VERSION}" \
    -o /go/bin/iac-gen ./cmd/iac-gen


# -------- prod stage -------- #
FROM alpine:3.17

WORKDIR /app/cafi-dev/iac-gen

# create non root user
RUN addgroup --gid 101 iac-gen && \
    adduser -S --uid 101 --ingroup iac-gen iac-gen

# run as non root user
USER 101

ENV PATH /go/bin:$PATH

# copy iac-gen binary from build
COPY --from=builder /go/bin/iac-gen bin/
COPY terraform ./terraform

EXPOSE 8000

ENTRYPOINT ["/app/cafi-dev/iac-gen/bin/iac-gen"]