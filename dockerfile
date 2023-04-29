FROM golang:1.20 as build

WORKDIR /app

COPY . .

RUN go build -o server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/server /server

EXPOSE 8080

ENTRYPOINT [ "/server" ]