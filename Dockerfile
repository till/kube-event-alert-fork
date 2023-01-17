FROM golang:1.19-alpine as build

WORKDIR "/build/kube-event-alert"

COPY config config
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN apk add --no-cache git
RUN go build -o kube-event-alert .

FROM alpine:3.17
RUN apk add ca-certificates
COPY --from=build /build/kube-event-alert/kube-event-alert /app/kube-event-alert

CMD "/app/kube-event-alert"
