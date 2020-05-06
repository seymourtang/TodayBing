FROM golang:latest
WORKDIR $GOPATH/src/TodayBing
COPY . $GOPATH/src/TodayBing
RUN GOPROXY="https://goproxy.io" GO111MODULE=on go build .
EXPOSE 5033
ENTRYPOINT ["./TodayBing"]
