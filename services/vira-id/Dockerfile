FROM golang:1.24.3-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# 1. Копируем только go.mod и go.sum — чтобы закэшировать deps
COPY go.mod go.sum ./

# 2. Загружаем зависимости
RUN go mod download

# 3. Теперь копируем остальной код (если выше ничего не изменилось — кеш сохранится)
COPY . .

# 4. Собираем бинарник
RUN go build -o vira-id .

# 5. Стартуем
CMD ["./vira-id"]
