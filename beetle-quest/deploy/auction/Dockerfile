# syntax=docker/dockerfile:1.7-labs
FROM golang:1.23 as build

WORKDIR /go/src/app
COPY --exclude=deploy --exclude=api . .

RUN go mod download
RUN CGO_ENABLED=0 go build  -o /go/bin/app ./cmd/auction/main.go

FROM gcr.io/distroless/base-debian12
COPY --from=build /go/bin/app /

# Set only in debug mode
ENV DEBUG=0

CMD ["/app"]
