FROM golang:1.18-alpine As build

# ENV GOPROXY=https://goproxy.cn,direct
RUN apk add --no-cache git
WORKDIR /app
ADD . /app
RUN cd /app && go build -o app

FROM alpine:latest
COPY --from=build /app/app /bin/
ARG traccarserveraddr ""
EXPOSE 8006
ENTRYPOINT /bin/app "-traccarserveraddr" ${traccarserveraddr}
