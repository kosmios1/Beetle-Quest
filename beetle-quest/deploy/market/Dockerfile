# syntax=docker/dockerfile:1.7-labs
FROM golang:1.23 as build

WORKDIR /go/src/app

COPY --exclude=deploy --exclude=api \
    --exclude=cmd/gacha/*   --exclude=internal/gacha/* \
    --exclude=cmd/user/*    --exclude=internal/user/* \
    --exclude=cmd/auth/*    --exclude=internal/auth/* \
    --exclude=cmd/admin/*   --exclude=internal/admin/* \
    --exclude=cmd/report/*  --exclude=internal/report/* \
    --exclude=cmd/static/* \
    ./go.mod ./go.mod

RUN go mod download
COPY --exclude=deploy --exclude=api . .

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG LD_FLAGS='-s -w'
ARG BUILD_FLAGS='-a'
ARG BUILD_TAGS='beetleQuestTest'

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} go build -tags=${BUILD_TAGS} -o /go/bin/app ${BUILD_FLAGS} -ldflags="${LD_FLAGS}" ./cmd/market/main.go


FROM gcr.io/distroless/base-debian12
COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/templates/ /templates/

EXPOSE 8080

CMD ["/app"]
