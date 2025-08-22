FROM golang:alpine as BUILDER
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV GOOS linux
ENV CGO_ENABLED 0
COPY . /usr/src/app/
RUN go build -o ./bin/ghwidgets ./cmd/server

FROM alpine:latest
RUN mkdir /app
ENV PORT 80
COPY --from=BUILDER /usr/src/app/bin/ghwidgets /app
CMD [ "/app/ghwidgets" ]
