FROM golang:1.17-alpine AS build
WORKDIR /temp/
COPY . ./
RUN go mod download
RUN go build -o /kubectl_top_exporter

FROM alpine:3.15
LABEL org.opencontainers.image.source https://github.com/jmnote/kubectl-top-exporter
WORKDIR /
COPY --from=build /kubectl_top_exporter /kubectl_top_exporter

EXPOSE     9977
USER       nobody
ENTRYPOINT ["/kubectl_top_exporter"]
