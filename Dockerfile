# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:latest as builder

# Manually install the Buffalo CLI
RUN go install -v github.com/gobuffalo/cli/cmd/buffalo@latest

ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /src/onlyoffice
WORKDIR /src/onlyoffice

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ADD . .
# Use the full path to the buffalo binary
RUN /go/bin/buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Run migrations first, THEN start the app.
# The '&&' means: "only start the app if the migrations succeed."
CMD /bin/app migrate && /bin/app