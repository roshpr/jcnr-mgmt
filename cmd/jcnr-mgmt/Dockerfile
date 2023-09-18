ARG BASE
FROM golang:1.21.0 as builder
WORKDIR /code
COPY *.go ./
COPY go.mod .
COPY go.sum .
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -trimpath -o /app .

FROM $BASE
COPY --from=builder /app .
