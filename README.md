# esia-mock

Mock-сервер для тестирования интеграции с ЕСИА (Единая Система Идентификации и Аутентификации).

## Возможности

- ✅ Полный OAuth2 flow (авторизация, выдача токенов)
- ✅ Веб-форма для ввода номера телефона
- ✅ **Уникальные моковые данные для каждого номера телефона** (in-memory кеш)
- ✅ Эндпоинты `/userinfo` и `/rs/prns/{oid}`
- ✅ Детальное логирование всех запросов

## Как работает кеширование данных

Когда пользователь вводит номер телефона на форме авторизации:

1. **Сервис сохраняет телефон** в in-memory кеше
2. **Генерирует уникальные моковые данные** на основе хеша телефона:
   - ФИО (из набора русских имен)
   - Дата рождения
   - Пол
   - СНILS (11 цифр)
   - ИНН (12 цифр)
   - Email
   - OID
   - Статусы (trusted, verified)

3. **При повторном входе** с тем же номером телефона возвращаются те же данные (детерминированная генерация)

### Пример

```bash
# Пользователь с телефоном +79991234567 всегда получит одни и те же данные:
{
  "oid": "1234567890",
  "firstName": "Иван",
  "lastName": "Иванов",
  "middleName": "Иванович",
  "birthDate": "15.03.1985",
  "gender": "M",
  "snils": "12345678901",
  "inn": "123456789012",
  "email": "ivan.ivanov.a1b2c3d@example.com",
  "mobile": "+79991234567",
  "trusted": true,
  "verified": true,
  "citizenship": "RUS",
  "status": "REGISTERED"
}

# Пользователь с телефоном +79109876543 получит другие данные
```

## Запуск

```bash
# Сборка
go build -o esia-mock ./cmd/app

# Запуск
./esia-mock
```

Сервер запустится на порту `8085`.

## Эндпоинты

- `GET /aas/oauth2/ac` - форма авторизации
- `GET /aas/oauth2/authorize` - форма авторизации
- `POST /aas/oauth2/authorize` - обработка формы авторизации
- `POST /aas/oauth2/te` - получение токена
- `GET /userinfo` - информация о пользователе (требует Bearer токен)
- `GET /rs/prns/{oid}` - информация о пользователе по OID (требует Bearer токен)

## Технические детали

### In-Memory кеш

- Хранилище реализовано в пакете `internal/storage`
- Thread-safe (использует `sync.RWMutex`)
- Генерация данных детерминированная (SHA256 от номера телефона)
- Данные живут до перезапуска сервера

### Структура проекта

```
.
├── cmd/app/
│   └── main.go              # Точка входа
├── internal/
│   ├── handler/
│   │   └── handler.go       # HTTP handlers
│   ├── logger/
│   │   └── logger.go        # Логирование
│   └── storage/
│       └── cache.go         # In-memory кеш с генерацией моковых данных
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Использование в тестах

```go
// 1. Настройте редирект на mock-сервер
clientID := "your-client-id"
redirectURI := "http://localhost:3000/callback"

// 2. Откройте форму авторизации
http://localhost:8085/aas/oauth2/authorize?client_id={clientID}&redirect_uri={redirectURI}&response_type=code&state=random_state

// 3. Введите любой номер телефона (например, +79991234567)

// 4. Получите код авторизации в redirect_uri

// 5. Обменяйте код на токен
POST http://localhost:8085/aas/oauth2/te
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code={code}&client_id={clientID}&redirect_uri={redirectURI}

// 6. Используйте access_token для получения данных пользователя
GET http://localhost:8085/userinfo
Authorization: Bearer {access_token}
```

## Зависимости

- `go.uber.org/zap` - структурированное логирование

## Лицензия

MIT
