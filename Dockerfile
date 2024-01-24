FROM golang:alpine as build-env

WORKDIR /app

COPY . ./

RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix  cgo -o  /example  ./cmd/eg/main.go


FROM alpine:latest

WORKDIR /

RUN mkdir /data
RUN addgroup --system example && adduser -S -s /bin/false -G expamle example

COPY --from=build-env  /example /example

RUN chown -R example:example /example
RUN chown -R example:example /data

USER example

EXPOSE 8080

ENTRYPOINT ["/example"]