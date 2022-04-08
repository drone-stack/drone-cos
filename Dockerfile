FROM ysicing/god AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/

COPY go.mod go.mod

COPY go.sum go.sum

RUN go mod download

COPY . .

WORKDIR /go/src/cmd

RUN go build -o ./cos

FROM ysicing/debian

COPY --from=builder /go/src/cmd/cos /bin/

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh /bin/cos

ENTRYPOINT ["/entrypoint.sh"]

CMD [ "/bin/cos" ]
