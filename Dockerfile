# Build the application from source
FROM golang:1.24.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /dnd-playlist-controller-server ./cmd/server

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /dnd-playlist-controller-server /dnd-playlist-controller-server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/dnd-playlist-controller-server"]

