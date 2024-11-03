FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
# Set environment variables (if needed)
ENV GOSUMDB=off
ENV GOPROXY=direct

# Install git
RUN apk update && apk add --no-cache git

WORKDIR /build

RUN go clean -modcache

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main main.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]
