FROM golang:1.19-alpine AS Builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /frolang

FROM alpine:latest
LABEL AUTHOR=akashsabu@ymail.com
LABEL VERSION=0.1.0
WORKDIR /
COPY --from=Builder /frolang /bin/
ENTRYPOINT ["frolang"]