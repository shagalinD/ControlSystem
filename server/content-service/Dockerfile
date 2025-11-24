FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Создаем директорию для загрузок
RUN mkdir -p /app/uploads

EXPOSE 8083

CMD ["./main"]