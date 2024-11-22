FROM golang:1.23 AS build

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build


FROM gcr.io/distroless/base-debian12

COPY --from=build /go/src/app/hermetic /hermetic
CMD ["/hermetic"]
