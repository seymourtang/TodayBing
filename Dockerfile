FROM golang:alpine AS build
RUN apk --no-cache add tzdata

RUN  mkdir /app
COPY . /app
WORKDIR /app
RUN GOPROXY="https://goproxy.io" CGO_ENABLED=0 GOOS=linux go build -o TodayBing

###
FROM scratch as fianl
COPY --from=build /app/TodayBing .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV GIN_MODE=release
ENV TZ=Asia/Shanghai
EXPOSE 5033

ENTRYPOINT ["./TodayBing"]
