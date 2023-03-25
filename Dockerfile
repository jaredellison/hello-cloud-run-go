# syntax=docker/dockerfile:1

## Build
FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Use CGO_ENABLED=0 to statically link binary so it can be used on the
# minimal distroless debian image
RUN CGO_ENABLED=0 go build -o ./bin/server cmd/server.go

## Deploy
# Using an image with only the bare essentials to run our app
# See docs on distroless images
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/bin/server /server

EXPOSE 8080

USER nonroot:nonroot

CMD [ "/server" ]
