FROM golang:1.19.6

WORKDIR /usr/src/shopware-metrics
COPY . /usr/src/shopware-metrics/

RUN go build -o /shopware-metrics .

CMD ["/shopware-metrics"]
