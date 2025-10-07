# 🏗️ Defect Manager - Система управления дефектами

## 📖 О проекте

**Менеджер дефектов** - это монолитное веб-приложение для централизованного управления дефектами на строительных объектах. Система обеспечивает полный цикл работы: от регистрации дефекта и назначения исполнителя до контроля статусов и формирования отчётности.

### 🎯 Основной функционал

- 🔐 **Аутентификация и авторизация** (JWT)
- 🏢 **Управление проектами** и строительными объектами
- 🐛 **Регистрация и отслеживание дефектов** с полным жизненным циклом
- 👥 **Назначение исполнителей** и контроль сроков
- 💬 **Комментарии и история изменений**
- 📎 **Прикрепление файлов** (фото, документы)
- 📊 **Аналитические отчеты** и экспорт данных
- 🔍 **Поиск, фильтрация и сортировка** дефектов

---

## 🛠 Технологический стек

### Backend

- **Go 1.21+** с фреймворком **Gin**
- **PostgreSQL** - основная база данных
- **GORM** - ORM для работы с БД
- **JWT** - аутентификация
- **bcrypt** - хеширование паролей

### Frontend

- **React 18** с **TypeScript**
- **Vite** - сборка и development server
- **React Router** - навигация
- **Tailwind CSS** - стилизация
- **Axios** - HTTP клиент
- **React Hook Form** - управление формами

---

## 🚀 Быстрый старт

### Предварительные требования

1. **Go** 1.21 или выше
2. **Node.js** 18.0 или выше
3. **PostgreSQL** 12 или выше
4. **Git**

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd kopatel_online
```

### 2. Настройка Backend

#### Установка зависимостей

```bash
cd server
go mod download
```

#### Настройка базы данных

1. Создайте базу данных в PostgreSQL:

```sql
CREATE DATABASE kopatel_online;
```

2. Создайте файл `.env` в папке `server/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_postgres_username
DB_PASSWORD=your_postgres_password
DB_NAME=kopatel_online
JWT_SECRET=your-super-secret-jwt-key-change-in-production
SERVER_PORT=8080
ENV=development
```

#### Запуск сервера

```bash
go run main.go
```

Сервер запустится на `http://localhost:8080`

### 3. Настройка Frontend

#### Установка зависимостей

```bash
cd client
npm install
```

#### Запуск development сервера

```bash
npm run dev
```

Frontend будет доступен на `http://localhost:5173`

---

## 📁 Структура проекта

```
kopatel_online/
├── server/                 # Backend приложение
│   ├── config/            # Конфигурация
│   ├── handlers/          # HTTP обработчики
│   ├── models/            # Модели данных
│   ├── postgres/          # Подключение к БД
│   ├── main.go            # Точка входа
│   └── .env               # Переменные окружения
├── client/                # Frontend приложение
│   ├── src/
│   │   ├── components/    # React компоненты
│   │   ├── pages/         # Страницы приложения
│   │   ├── hooks/         # Кастомные хуки
│   │   ├── services/      # API клиент
│   │   ├── types/         # TypeScript типы
│   │   └── utils/         # Вспомогательные функции
│   ├── public/            # Статические файлы
│   └── package.json
└── README.md
```

---

## 👥 Роли пользователей

### 1. 👷 Инженер (Engineer)

- Регистрация новых дефектов
- Обновление информации по дефектам
- Добавление комментариев и файлов
- Просмотр назначенных дефектов

### 2. 👨‍💼 Менеджер (Manager)

- Управление проектами
- Назначение исполнителей на дефекты
- Контроль сроков выполнения
- Формирование отчетов

### 3. 👀 Наблюдатель (Observer)

- Просмотр проектов и дефектов
- Мониторинг прогресса
- Просмотр отчетности

---

## 🔑 Начальная настройка

### Создание первого пользователя

После запуска сервера автоматически создаются роли. Для создания первого пользователя отправьте POST запрос:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "admin123",
    "full_name": "Администратор",
    "role_id": 2
  }'
```

### Доступные роли при первом запуске:

- `1` - Инженер (engineer)
- `2` - Менеджер (manager)
- `3` - Наблюдатель (observer)

---

## 📚 API Endpoints

### Аутентификация

- `POST /auth/register` - Регистрация
- `POST /auth/login` - Вход в систему

### Проекты

- `GET /api/projects` - Список проектов
- `POST /api/projects` - Создание проекта
- `GET /api/projects/:id` - Получение проекта
- `PUT /api/projects/:id` - Обновление проекта
- `DELETE /api/projects/:id` - Удаление проекта

### Дефекты

- `GET /api/defects` - Список дефектов (с фильтрацией)
- `POST /api/defects` - Создание дефекта
- `GET /api/defects/:id` - Получение дефекта
- `PUT /api/defects/:id` - Обновление дефекта
- `PATCH /api/defects/:id/status` - Изменение статуса
- `DELETE /api/defects/:id` - Удаление дефекта
- `GET /api/defects/my` - Мои дефекты

### Комментарии

- `GET /api/comments/defect/:defect_id` - Комментарии дефекта
- `POST /api/comments/defect/:defect_id` - Добавление комментария

### Вложения

- `POST /api/attachments/defect/:defect_id` - Загрузка файла
- `GET /api/attachments/defect/:defect_id` - Список файлов
- `GET /api/attachments/:id/download` - Скачивание файла

### Отчеты

- `GET /api/reports/defects` - Аналитика по дефектам
- `GET /api/reports/defects/export` - Экспорт в CSV
- `GET /api/reports/projects/:project_id` - Отчет по проекту

---

## 🐛 Жизненный цикл дефекта

1. **Новая** → Создан, ожидает назначения
2. **В работе** → Назначен исполнитель, работа ведется
3. **На проверке** → Работы завершены, ожидает проверки
4. **Закрыта** → Дефект устранен
5. **Отменена** → Дефект не требует устранения

---

## 🛠 Команды разработки

### Backend команды

```bash
cd server

# Запуск в development режиме
go run main.go

# Запуск тестов
go test ./handlers/...

# Сборка для production
go build -o kopatel-server main.go
```

### Frontend команды

```bash
cd client

# Development сервер
npm run dev

# Production сборка
npm run build

# Preview production сборки
npm run preview

# Запуск тестов
npm test
```

---

## 🔧 Конфигурация

### Backend переменные окружения (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=kopatel_online
JWT_SECRET=very-secret-key
SERVER_PORT=8080
ENV=development
```

## 📊 Особенности системы

### Производительность

- Время отклика страницы ≤ 1 секунды
- Поддержка 50+ одновременных пользователей
- Оптимизированные запросы к БД

### Безопасность

- 🔒 JWT аутентификация
- 🔐 Хеширование паролей (bcrypt)
- 🛡️ Защита от SQL-инъекций
- 🚫 Защита от XSS и CSRF атак
- 📝 Логирование действий пользователей

### Совместимость

- **Браузеры:** Chrome, Firefox, Edge (последние версии)
- **Платформы:** ПК, планшеты
- **Язык:** Русский

---

## 🚀 Деплой в Production

### Backend

1. Соберите бинарный файл: `go build -o kopatel-server main.go`
2. Настройте reverse proxy (Nginx)
3. Настройте SSL сертификат
4. Запустите как systemd service

### Frontend

1. Соберите проект: `npm run build`
2. Разместите файлы из `dist/` на веб-сервере
3. Настройте маршрутизацию для SPA

### База данных

1. Настройте регулярное резервное копирование
2. Настройте репликацию (опционально)
3. Оптимизируйте индексы

---
