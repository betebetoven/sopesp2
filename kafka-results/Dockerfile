FROM golang:1.21-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o kafka-results

EXPOSE 8082

CMD ["./kafka-results"]