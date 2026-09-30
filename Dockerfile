# syntax=docker/dockerfile:1

FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/contract ./cmd/contract
RUN go run ./cmd/openapi-bundle api/contract.yaml /out/build/openapi.bundle.yaml

FROM gcr.io/distroless/static-debian12:nonroot
# GET /api/v2/contracts/openapi reads build/openapi.bundle.yaml relative to WORKDIR
WORKDIR /app
COPY --from=build /out/contract /app/contract
COPY --from=build /out/build/openapi.bundle.yaml /app/build/openapi.bundle.yaml
EXPOSE 8082
USER nonroot:nonroot
ENTRYPOINT ["/app/contract"]
