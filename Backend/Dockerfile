# Development Stage
FROM golang:1.24-alpine

# install git (needed by go modules) and other tools
RUN apk add --no-cache git

# set env for static builds
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# install Air for live-reload
RUN go install github.com/air-verse/air@latest

# cache deps
COPY Backend/go.mod Backend/go.sum ./
RUN go mod download

# copy source code
COPY Backend/ ./

# (optional) expose your app port; adjust if you listen on another port
EXPOSE 8080

# run Air on container start
CMD ["air"]