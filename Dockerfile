FROM golang:1.20-alpine AS builder

# Set necessary environmet variables needed for our image
ENV CGO_ENABLED=0

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
COPY client.go .
COPY module.go .
COPY stats.go .

RUN go install go.k6.io/xk6/cmd/xk6@latest
RUN xk6 build \
    --with github.com/gjergjsheldija/xk6-mllp=.

#############################
## STEP 2 build a small image
#############################
FROM scratch

COPY --from=builder /build/k6 /

# Command to run the executable
ENTRYPOINT ["/k6"]
