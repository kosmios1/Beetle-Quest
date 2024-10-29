# syntax=docker/dockerfile:1.7-labs
FROM golang:1.23 as build

WORKDIR /go/src/app

COPY --exclude=deploy --exclude=api ./go.mod ./go.mod
RUN go mod download

COPY --exclude=deploy --exclude=api . .
RUN CGO_ENABLED=0 go build  -o /go/bin/app ./cmd/user/main.go

FROM gcr.io/distroless/base-debian12
COPY --from=build /go/bin/app /

EXPOSE 8080

CMD ["/app"]
