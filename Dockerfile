FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
# Set environment variables (if needed)
ENV GOSUMDB=off
ENV GOPROXY=direct

# Install git
RUN apk update && apk add --no-cache git && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /build

#RUN go clean -modcache

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main main.go


FROM scratch

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /build/.env /app/.env
COPY --from=builder /build/root.crt /app/.env

EXPOSE 8098

CMD ["./main"]
