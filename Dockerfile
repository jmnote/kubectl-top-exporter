FROM golang:1.17-alpine AS build
WORKDIR /temp/
COPY . ./
RUN go mod download
RUN go build -o /kubectl-top-exporter

FROM alpine:3.15
WORKDIR /
COPY --from=build /kubectl-top-exporter /kubectl-top-exporter
EXPOSE 9100
ENTRYPOINT ["/kubectl-top-exporter"]
