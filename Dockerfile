FROM --platform=$BUILDPLATFORM golang:1.17.6-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /src ./cmd/flareup

FROM alpine
COPY --from=build /src/flareup /bin
CMD ["/bin/flareup"]