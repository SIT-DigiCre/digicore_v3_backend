FROM golang:1.17.3 as development

WORKDIR "/app"
COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/cosmtrek/air

COPY . .

CMD ["air"]


FROM golang:1.17.3 as build

WORKDIR /work
COPY . .
RUN go mod download
RUN go build


FROM gcr.io/distroless/base-debian10 as production

COPY --from=build /work/digicore_v3_backend /
COPY client_secret.json /

CMD ["/digicore_v3_backend"]
