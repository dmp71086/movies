# MULTI STAGE EXAMPLE

#######################################
# STAGE 1. BUILD STAGE
#######################################

# Для сборки нашего приложения на go требуется образ ОС, в котором установлен golang нужной нам версии.
# Alpine выбран из-за его небольшого размера по сравнению с Ubuntu.
FROM golang:1.22-alpine3.19 AS build

# Устанавливаем место назначения для COPY
WORKDIR /app

# Копируем файлы go.mod и go.sum в WORKDIR
COPY go.* ./
# Скачиваем необходимые Go модули (зависимости нашего проетка)
RUN go mod download

# Копируем все исходные go файлы нашего проекта в образ
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Собираем бинарный файл нашего приложения
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/movies_server ./cmd

#######################################
# STAGE 2. FINAL STAGE
#######################################

FROM scratch AS final

WORKDIR /

COPY --from=build /bin/movies_server /movies_server
COPY --from=build /app/cert /cert

# Указываем какой порт необходимо слушать
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8888 8888
EXPOSE 8000 8000

# Точка входа
ENTRYPOINT ["/movies_server"]