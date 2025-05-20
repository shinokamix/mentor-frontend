# MentorLink 

Платформа для поиска менторов. Микросервисная архитектура с авторизацией, отзывами, рейтингом и централизованным логгированием.

---

## Технологии

- **Go** — основной язык
- **Kafka + Kafka UI** — шина событий
- **PostgreSQL + golang-migrate** — база и миграции
- **Redis** — быстрая память
- **Docker / Docker Compose** — всё в контейнерах
- **gRPC** — общение между сервисами
- **Promtail + Loki + Grafana** — логгирование
- **nginx** — API Gateway

---

## Архитектура

> Схема текущей архитектуры проекта:

![architecture](./monitoring/architecture.jpg)

---

##  Быстрый старт

> Запуск — миграции, сервисы и логгирование.

```bash
# 1. Клонируем репозиторий
git clone https://github.com/GkadyrG/MentorLink.git
cd MentorLink

# 2. Запускаем всё, пересобирая образы
docker compose up -d --build

# 3. Применяем миграции (внутри хоста)
make migrate-up-all