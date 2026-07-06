# --- build stage ---
FROM golang:1.23-alpine AS build

WORKDIR /src

# Cache module downloads.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary; templates and JS are embedded via //go:embed.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

# --- runtime stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=build /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
