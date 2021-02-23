
FROM golang:1.15.2-alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
COPY mllp.go .
COPY client.go .
RUN go mod download

RUN go get -u github.com/k6io/xk6/cmd/xk6
RUN xk6 build --with github.com/gjergjsheldija/xk6-mllp=. 

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /build/k6 /

# Command to run the executable
ENTRYPOINT ["/k6"]