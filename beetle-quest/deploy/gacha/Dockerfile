# syntax=docker/dockerfile:1.7-labs
FROM golang:1.23 as build

WORKDIR /go/src/app


COPY --exclude=deploy --exclude=api \
    --exclude=cmd/market/*  --exclude=internal/market/* \
    --exclude=cmd/user/*    --exclude=internal/user/* \
    --exclude=cmd/auth/*    --exclude=internal/auth/* \
    --exclude=cmd/admin/*   --exclude=internal/admin/* \
    --exclude=cmd/report/*  --exclude=internal/report/* \
    --exclude=cmd/static/* \
    ./go.mod ./go.mod

RUN go mod download

COPY --exclude=deploy --exclude=api . .
RUN CGO_ENABLED=0 go build  -o /gacha-service ./cmd/gacha/main.go

FROM gcr.io/distroless/base-debian12
COPY --from=build /gacha-service /
COPY --from=build /go/src/app/templates/ /templates/

EXPOSE 8080

CMD ["/gacha-service"]
