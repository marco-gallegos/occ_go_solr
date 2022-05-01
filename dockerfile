FROM golang:alpine AS build
WORKDIR /go/src/myapp
COPY . .
RUN go build -o /go/bin/myapp main.go

FROM scratch
COPY --from=build /go/bin/myapp /go/bin/myapp
EXPOSE 8080
ENTRYPOINT ["/go/bin/myapp"]