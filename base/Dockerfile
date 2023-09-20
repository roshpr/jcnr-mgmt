#FROM gcr.io/distroless/base
FROM docker.io/library/ubuntu-go:latest
# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single
CMD ["./app"]
