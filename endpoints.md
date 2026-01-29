
**Часть 12: Документация API `clofi` — полная спецификация эндпоинтов**

Ниже приведена полная и точная документация всех REST-эндпоинтов вашего сервиса, включая методы, параметры, тела запросов, ответы, коды ошибок и примеры `curl`.

---

## 🔐 Общие соглашения

- **Базовый URL**: `http://localhost:8080`
- **Формат запросов/ответов**: `application/json`
- **Авторизация**: заголовок `Authorization: Bearer <JWT>`
- **Кодировка**: UTF-8

---

## 🧾 1. Регистрация пользователя

### HTTP метод и путь
```
POST /register
```

### Параметры запроса
— Нет query/path параметров  
— Нет заголовков (кроме `Content-Type: application/json`)

### Тело запроса (JSON)
```json
{
  "username": "string (3–32 символа, только буквы/цифры)",
  "name": "string (1–64 символа)",
  "password": "string (минимум 6 символов)"
}
```

### Формат успешного ответа
```json
{
  "message": "пользователь создан"
}
```

### Коды состояния
| Код | Описание |
|-----|--------|
| `201 Created` | Пользователь успешно создан |
| `400 Bad Request` | Неверный JSON или нарушение валидации |
| `409 Conflict` | Пользователь с таким логином уже существует |

### Пример curl
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "name": "Alice Cooper",
    "password": "secure123"
  }'
```

---

## 🔑 2. Вход в систему (получение JWT)

### HTTP метод и путь
```
POST /login
```

### Параметры запроса
— Нет query/path параметров

### Тело запроса (JSON)
```json
{
  "username": "string",
  "password": "string"
}
```

### Формат успешного ответа
```json
{
  "token": "string (JWT)"
}
```

### Коды состояния
| Код | Описание |
|-----|--------|
| `200 OK` | Успешный вход |
| `400 Bad Request` | Неверный JSON |
| `401 Unauthorized` | Неверный логин/пароль |

### Пример curl
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "secure123"
  }'
```

---

## 🔍 3. Получение списка товаров

### HTTP метод и путь
```
GET /products
```

### Параметры запроса (query)
| Параметр | Тип | Обязательный | Описание |
|--------|-----|------------|---------|
| `page` | int | нет | Номер страницы (по умолчанию: 1) |
| `limit` | int | нет | Товаров на страницу (по умолчанию: 20, максимум: 100) |
| `q` | string | нет | Текст для поиска (по `name` и `description`) |
| `brand` | string | нет | Фильтр по бренду (точное совпадение) |
| `availability` | bool | нет | `true` — только в наличии, `false` — только отсутствующие |
| `sort_price` | string | нет | `asc` — по возрастанию цены, `desc` — по убыванию |

### Тело запроса
— Отсутствует

### Формат успешного ответа
Массив товаров:
```json
[
  {
    "id": "string (MongoDB ObjectId)",
    "name": "string",
    "url": "string (URL)",
    "price": "int (в рублях)",
    "brand": "string",
    "photos_urls": ["string (URL)", ...],
    "availability": "bool",
    "description": "string"
  }
]
```

### Коды состояния
| Код | Описание |
|-----|--------|
| `200 OK` | Успешно |
| `500 Internal Server Error` | Ошибка базы данных |

### Примеры curl

#### Все товары (первая страница):
```bash
curl "http://localhost:8080/products"
```

#### Поиск по тексту:
```bash
curl "http://localhost:8080/products?q=джинсы"
```

#### Фильтрация + сортировка:
```bash
curl "http://localhost:8080/products?brand=Levis&availability=true&sort_price=asc&page=1&limit=10"
```

---

## 📦 4. Получение товара по ID

### HTTP метод и путь
```
GET /product
```

### Параметры запроса (query)
| Параметр | Тип | Обязательный | Описание |
|--------|-----|------------|---------|
| `id` | string | **да** | ID товара (MongoDB ObjectId) |

### Тело запроса
— Отсутствует

### Формат успешного ответа
```json
{
  "id": "string",
  "name": "string",
  "url": "string",
  "price": "int",
  "brand": "string",
  "photos_urls": ["string", ...],
  "availability": "bool",
  "description": "string"
}
```

### Коды состояния
| Код | Описание |
|-----|--------|
| `200 OK` | Товар найден |
| `400 Bad Request` | Не указан `id` |
| `404 Not Found` | Товар не найден |
| `500 Internal Server Error` | Ошибка БД |

### Пример curl
```bash
curl "http://localhost:8080/product?id=65a1b2c3d4e5f67890123456"
```

---

## ❤️ 5. Переключение лайка

### HTTP метод и путь
```
POST /like
```

### Параметры запроса (query)
| Параметр | Тип | Обязательный | Описание |
|--------|-----|------------|---------|
| `product_id` | string | **да** | ID товара |

### Заголовки
- `Authorization: Bearer <JWT>` (**обязательно**)

### Тело запроса
— Отсутствует

### Формат успешного ответа
```json
{
  "status": "лайк добавлен"  // или "лайк удалён"
}
```

### Коды состояния
| Код | Описание |
|-----|--------|
| `200 OK` | Лайк переключён |
| `400 Bad Request` | Не указан `product_id` |
| `401 Unauthorized` | Отсутствует/недействителен JWT |
| `404 Not Found` | Пользователь или товар не найден |
| `500 Internal Server Error` | Ошибка БД |

### Пример curl
```bash
curl -X POST "http://localhost:8080/like?product_id=65a1b2c3d4e5f67890123456" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 📌 Общие ошибки (для всех эндпоинтов)

Если возникает ошибка, ответ всегда имеет формат:
```json
{
  "error": "описание ошибки на русском"
}
```

Пример:
```json
{ "error": "неверный логин или пароль" }
```

---

## ✅ Итог

| Эндпоинт | Метод | Авторизация | Назначение |
|---------|------|------------|-----------|
| `/register` | POST | ❌ | Регистрация |
| `/login` | POST | ❌ | Получение JWT |
| `/products` | GET | ❌ | Поиск и фильтрация товаров |
| `/product` | GET | ❌ | Получение товара по ID |
| `/like` | POST | ✅ | Поставить/убрать лайк |

Этот API полностью готов к использованию фронтендом, мобильным приложением или другим клиентом.
