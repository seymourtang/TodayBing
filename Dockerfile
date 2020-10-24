FROM golang:alpine AS build
RUN apk --no-cache add tzdata

RUN  mkdir /app
COPY . /app
WORKDIR /app
RUN  CGO_ENABLED=0 GOOS=linux go build -o todaybing ./cmd/

###
FROM scratch as fianl
COPY --from=build /app/todaybing .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
EXPOSE 5033

ENTRYPOINT ["./todaybing"]
