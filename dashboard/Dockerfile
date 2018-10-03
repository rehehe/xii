FROM golang as builder
LABEL stage="intermediate"

COPY . /dashboard
WORKDIR /dashboard

RUN go get -v \
&&  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o build/linux/dashboard cmd/dashboard/main.go

###
FROM scratch as stag
LABEL description="xii.dashboard"

COPY --from=builder /dashboard/build/linux/dashboard /

EXPOSE 1101

ENTRYPOINT ["./dashboard"]


