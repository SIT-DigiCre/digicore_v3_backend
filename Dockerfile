FROM golang:1.20.0 as build

WORKDIR /work
COPY . .
RUN go mod download
RUN go build ./cmd/digicore_v3_backend

FROM gcr.io/distroless/base-debian10 as production

COPY --from=build /work/digicore_v3_backend /

CMD ["/digicore_v3_backend"]

FROM golang:1.20.0 as admin

WORKDIR "/app"
RUN go install github.com/k0kubun/sqldef/cmd/mysqldef@v0.13.19
RUN apt-get update && apt-get install -y default-mysql-client-core

COPY . .
