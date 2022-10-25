FROM golang:1.19-alpine AS build
WORKDIR /app
COPY . ./
RUN go build -o /frolang

FROM alpine:latest
WORKDIR /
COPY --from=build /frolang /frolang
ENTRYPOINT ["/frolang"]