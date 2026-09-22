FROM golang:2.25.4-bookworm
WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

CMD ["",""]
