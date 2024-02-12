FROM golang:1.22 AS build-stage

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /hermetic

FROM gcr.io/distroless/base-debian11 AS run-stage

WORKDIR /run

COPY --from=build-stage /hermetic .

CMD ["./hermetic"]
