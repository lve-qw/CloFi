## Требования

    ОС: Linux (Ubuntu 22.04+ рекомендуется)
    Docker 
    Docker Compose
    Go 1.22+

---

## Для запустка 
    git clone https://github.com/yourname/CloFi.git
    
    cd CloFi

    docker-compose up --build



---

## Конфигурация

Переменные окружения (задать в `.env`):

| Переменная | По умолчанию | Описание |
|-----------|-------------|---------|
| `POSTGRES_HOST` | `postgres` | Хост PostgreSQL |
| `MONGO_URI` | `mongodb://mongo:27017` | URI MongoDB |
| `JWT_SECRET` | `my_super_secret_key_123!` | Секрет для JWT |
| `SERVER_PORT` | `8080` | Порт сервера |

---
