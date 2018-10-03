FROM golang as builder
LABEL stage="intermediate"

COPY . /survey
WORKDIR /survey

RUN go get -v \
&&  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o build/linux/survey cmd/survey/main.go

###
FROM scratch as stag
LABEL description="xii.survey"

COPY --from=builder /survey/build/linux/survey /

EXPOSE 1101

ENTRYPOINT ["./survey"]
