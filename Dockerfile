FROM golang:1.19.0 as development

WORKDIR "/app"
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY . .

CMD ["air", "-c", "./config/.air.toml"]


FROM golang:1.19.0 as build

WORKDIR /work
COPY . .
RUN go mod download
RUN go build ./cmd/digicore_v3_backend


FROM gcr.io/distroless/base-debian10 as production

COPY --from=build /work/digicore_v3_backend /

CMD ["/digicore_v3_backend"]


FROM golang:1.19.0 as admin

WORKDIR "/app"
RUN apt-get update && apt-get install -y default-mysql-client-core

COPY . .
