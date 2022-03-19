FROM golang:1.17.3

WORKDIR "/app"
COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/cosmtrek/air

COPY . .

CMD ["air"]
