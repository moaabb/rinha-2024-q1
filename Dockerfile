# Build the application from source
FROM golang:1.22.0-alpine3.19 AS build-stage

WORKDIR /app

ADD . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o rinha-go ./cmd/*.go

# Run the tests in the container
FROM build-stage AS run-test-stage

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/rinha-go /app/rinha-go

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/rinha-go"]