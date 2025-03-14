# Распределенный калькулятор арифметических выражений  
  
## Описание  
Это проект распределенного калькулятора, который обрабатывает арифметические выражения. Система состоит из двух основных компонентов: оркестратора и агентов. Служба оркестратора управляет входящими запросами на вычисления и распределяет рабочую нагрузку между несколькими экземплярами агентов. Агенты выполняют фактические арифметические вычисления с настраиваемым временем обработки для различных операций (сложение, вычитание, умножение, деление).  
  
## Возможности  
  
Калькулятор обладает следующими возможностями:  
  
1. Арифметические операции: Базовые операции: сложение ( + ), вычитание ( - ), умножение ( * ), деление ( / ), обработка десятичных чисел с высокой точностью, Поддержка очень больших и очень маленьких чисел, правильная обработка приоритета операторов.  
2. Функции выражений: Поддержка скобок для вложенных выражений: (2 + 3) * (4 + 5), унарный оператор минус в разных контекстах (-2, 2 * -3), несколько операций в одном выражении, сложные вложенные выражения, гибкая обработка пробелов.  
3. Проверка ввода: Проверка пустых выражений, проверка сбалансированных скобок, проверка использования десятичной точки, предотвращение недопустимых символов, проверка последовательных операторов, проверка отсутствующих операндов/операторов, защита от деления на ноль.  
4. Распределенная обработка: Параллельная обработка вычислений, распределение задач по нескольким агенты, конфигурация времени работы для различных операций, ведение журнала запросов/ответов, обработка ошибок и отслеживание статуса.  
5. Дополнительные функции: Отслеживание статуса выражения (ожидание, в процессе, завершено, ошибка), подробный отчет об ошибках, комплексная система журналирования, поддержка длинных выражений, высокоточные десятичные вычисления.  
  
## Предварительные требования  
  
Перед началом установки убедитесь, что у вас установлены:  
  
1. **Go версии 1.24.0**  
- Скачайте Go с [официального сайта](https://golang.org/dl/)  
- Проверьте установку: `go version`  
  
2. **Git**  
- Установите Git с [git-scm.com](https://git-scm.com/)  
- Проверьте установку: `git --version`  
  
3. **Docker** (опционально, для запуска в контейнерах)  
- Установите Docker с [docker.com](https://www.docker.com/)  
- Проверьте установку: `docker --version`  
  
4. **Make** (опционально, для использования Makefile)  
- Windows: установите через [Chocolatey](https://chocolatey.org/): `choco install make`  
- Проверьте установку: `make --version`  
  
## Использованные технологии  
  
- **Go v1.24.0** - Основной язык программирования  
- **Стандартные библиотеки Go** - Базовая функциональность  
- **Zap logger v1.27.0** - Высокопроизводительное логирование  
- **Google генератор UUID v1.6.0** - Генерация уникальных идентификаторов  
- **Testify для улучшения тестирования v1.10.0** - Улучшенное тестирование  
- **Роутер HTTP gorilla/mux v1.8.1** - Маршрутизация HTTP запросов  
- **Отладочный вывод go-spew v1.1.1** - Улучшенная отладка  
- **Для вычисления и форматирования разницы между строками - go-difflib v1.0.0** - Сравнение строк  
- **Библиотека для объединения нескольких ошибок в одну - multierr v1.11.0** - Обработка ошибок  
- **Библиотека для работы с YAML - yaml.v3 v3.0.1** - Работа с конфигурациями  
  
## Установка  
  
### Клонируйте репозиторий  
  
**Используя HTTPS:**  
```  
git clone https://github.com/flexer2006/y.lms-sprint2-calculator.git  
```  
  
**Или используя SSH:**  
```  
git clone git@github.com:flexer2006/y.lms-sprint2-calculator.git  
```  
  
## Запуск   
  
### Перейдите в каталог проекта:  
```  
cd y.lms-sprint2-calculator  
```  
  
### Метод 1: Использование Docker  
  
#### 1. Сборка образов  
  
##### Сборка для конкретной платформы:  
  
**Сборка отдельных образов 
  
```bash
 docker-compose up -d --build
```
```bash
docker-compose up -d --scale agent=N
```
N -количество агентов

**Остановка:**

```bash
docker-compose stop
```

**Удаление:**


```bash
docker-compose down
```

### 1.5. Метод: Использовать докер команды

**Сборка отдельных образов по целевым этапам**

```bash
docker build --target orchestrator -t my-app-orchestrator .
```

```bash
docker build --target agent -t my-app-agent .
```

**Запустите оркестратор:**  
```bash 
docker run -d -p 8080:8080 --name orchestrator my-app-orchestrator  
```  
Флаги:  
- `-d`: запуск в фоновом режиме  
- `-p 8080:8080`: проброс порта 8080 из контейнера на хост  
- `--name orchestrator`: имя контейнера  
**С использованием сети (рекомендуемый способ):**  
```bash 
docker network create calculator-network  
docker run -d -p 8080:8080 --name orchestrator --network calculator-network calculator-orchestrator  
```  
  
**Запустите одного или несколько агентов:**  
```bash
docker run -d --name agent1 my-app-agent  
docker run -d --name agent2 my-app-agent  
docker run -d --name agent3 my-app-agent  
```  
Вы можете запустить столько агентов, сколько необходимо для обработки нагрузки.  
  
#### Управление контейнерами  
  
**Просмотр запущенных контейнеров:**  
```bash 
docker ps  
```  
  
**Остановка контейнеров:**  
```bash
docker stop orchestrator agent1 agent2 agent3  
```  
  
**Удаление контейнеров:**  
```bash  
docker rm orchestrator agent1 agent2 agent3  
```  
  
### Метод 2: Использование Makefile  
  
Makefile предоставляет удобные команды для сборки и запуска проекта.  
  
#### Доступные команды:  
  
**Просмотр всех доступных команд:**  
```bash 
make help  
```  
  
**Сборка проекта:**  
```bash 
make build  
```  
Соберет оба сервиса (оркестратор и агент)  
  
**Если хотите включить сборку с gcc (CGO_ENABLED=1), выполните:**  
```bash  
ENABLE_CGO=1 make build  
```  
  
**Запуск в различных режимах:**  
  
- **Стандартный запуск:**  
```bash 
make run  
```  
Запускает оркестратор и одного агента с настройками по умолчанию  
  
- **Запуск в режиме разработки:**  
```bash 
make run-dev  
```  
Запускает с расширенным логированием и дополнительной отладочной информацией  
  
- **Запуск в производственном режиме:**  
```bash  
make run-prod  
```  
Запускает с оптимизированными настройками для производственной среды  
  
- **Запуск только агента:**  
```bash 
make run-agent  
```  
  
- **Запуск только оркестратора:**  
```bash 
make run-orchestrator  
```  
  
- **Запуск с проверкой гонок:**  
```bash 
make run-race  
```  
- **Запуск большего количества агентов:**  
  
```bash  
COMPUTING_POWER=4 make run-agent  
```  
  
### Метод 3: Использование PowerShell  
  
PowerShell скрипты предоставляют простой способ сборки и запуска проекта в Windows.  
  
#### 1. Сборка проекта  
  
**Запустите скрипт сборки:**  
```bash  
.\build.ps1 build  
```  
Этот скрипт выполнит:  
- Проверку зависимостей  
- Компиляцию оркестратора и агента  
- Подготовку конфигурационных файлов  
  
#### 2. Запуск в различных режимах  
  
**Запуск в режиме разработки:**  
```bash  
.\build.ps1 run-dev  
```  
Запускает с расширенным логированием  
  
**Запуск в производственном режиме:**  
```bash 
.\build.ps1 run-prod  
```  
Запускает с оптимизированными настройками  
  
**Остановить:**  
```bash  
.\build.ps1 stop  
```  
**Запуск с большм количеством агентов:**  
```bash  
$env:COMPUTING_POWER = 4  
.\build.ps1 run-agent  
```  
  
## Переменные среды  
  
### Служба оркестратора:  
  
1. `PORT` - HTTP порт сервера (по умолчанию: 8080)  
2. `TIME_ADDITION_MS` - Время выполнения операции сложения (по умолчанию: 100)  
3. `TIME_SUBTRACTION_MS` - Время выполнения операции вычитания (по умолчанию: 100)  
4. `TIME_MULTIPLY_MS` - Время выполнения операции умножения (по умолчанию: 200)  
5. `TIME_DIVISION_MS` - Время выполнения операции деления (по умолчанию: 200)  
  
### Агентская служба:  
  
1. `COMPUTING_POWER` - Количество одновременных вычислений (по умолчанию: 1)  
2. `ORCHESTRATOR_URL` - URL-адрес службы оркестратора (по умолчанию: http://localhost:8080)  
  
## API эндпоинты  
  
### Служба оркестратора:  
  
- `POST /api/v1/calculate` - Отправить выражение для вычисления  
- `GET /api/v1/expressions` - Список всех выражений  
- `GET /api/v1/expressions/{id}` - Получить статус и результат выражения  
  
### Внутренний API  
  
- `GET /internal/task` - Получить следующую задачу (используется агентами)  
- `POST /internal/task` - Отправить результат задачи (используется агентами)  
  
## Примеры использования  
  
### Успешное вычисление  
  
**Используя curl запрос:**  
```  
curl -L 'http://localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' --data '{"expression":"2+2*2"}'  
```  
**Ответ (HTTP 201 Created):**  
```json  
  
{  
  
    "id": "123e4567-e89b-12d3-a456-426614174000"  
}  
  
```  
**Проверка решения примера:**  
```  
curl --location 'http://localhost:8080/api/v1/expressions/123e4567-e89b-12d3-a456-426614174000'  
```  
**Ответ (HTTP 200 OK):**  
```json  
  
{  
  
    "expression": {  
        "id": "123e4567-e89b-12d3-a456-426614174000",  
        "expression": "2+2*2",  
        "status": "COMPLETE",  
        "result": 6  
    }  
}  
  
```  
  
### Неуспешное вычисление  
  
**Пустое выражение:**  
```  
curl -L 'http://localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' --data '{"expression":""}'  
```  
**Ответ (HTTP 422 Unprocessable Entity):**  
```json  
  
{  
  
    "error": "Expression cannot be empty"  
}  
  
```  
  
**Неверный формат выражения:**  
```  
curl -L 'http://localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' --data '{"expression":"2++2"}'  
```  
**Ответ (HTTP 201 Created, но проверьте статус):**  
```json  
  
{  
  
    "expression": {  
        "id": "123e4567-e89b-12d3-a456-426614174000",  
        "expression": "2++2",  
        "status": "ERROR",  
        "error": "invalid expression: too few tokens"  
    }  
}  
  
```  
  
**Выражение не найдено:**  
```bash  
curl --location 'http://localhost:8080/api/v1/expressions/non-existent-id'  
```  
Ответ (HTTP 404 Not Found):  
```json  
  
{  
  
    "error": "Expression not found"  
}  
  
```  
  
### Список всех выражений  
  
**Ввод:**  
```bash  
curl --location 'http://localhost:8080/api/v1/expressions'  
```  
**Ответ (HTTP 200 OK):**  
```json  
  
{  
  
    "expressions": [  
        {  
            "id": "123e4567-e89b-12d3-a456-426614174000",  
            "expression": "2+2*2",  
            "status": "COMPLETE",  
            "result": 6  
        },  
        {  
            "id": "987fcdeb-51d3-12a4-b678-426614174000",  
            "expression": "10-5",  
            "status": "COMPLETE",  
            "result": 5  
        }  
    ]  
}  
  
```  
  
### Postman  
  
Вы можете использовать Postman для проверки проекта:  
![Pasted image 20250221114952](https://github.com/user-attachments/assets/a371d09d-d783-472e-af80-e1898432c25e)  
  
  
### Структура проекта  
  
Проект организован следующим образом:  
- `cmd/` - Точки входа для оркестратора и агента  
- `internal/` - Внутренняя логика приложения  
- `pkg/` - Переиспользуемые пакеты  
- `tests/` - Тестовые файлы  
- `configs/` - Конфигурационные файлы  
  
### Тестирование  
  
**Запуск всех тестов:**  
```bash  
go test ./...  
```  
  
**Подробный вывод тестов:**  
```bash 
go test ./... -v  
```
  
**Запуск конкретного теста:**  
```bash 
go test ./tests -run TestCalculator  
```
  
## Устранение неполадок  
  
### Общие проблемы  
  
1. **Ошибка "connection refused"**  
- Убедитесь, что оркестратор запущен и доступен  
- Проверьте правильность URL в конфигурации агента  
  
2. **Ошибка компиляции**  
- Выполните `go mod tidy` для обновления зависимостей  
- Проверьте версию Go (требуется 1.23.6 или выше)  
  
3. **Агент не подключается к оркестратору**  
- Проверьте сетевые настройки  
- Убедитесь, что порты не заблокированы файрволом  
  
## Схема  
![image](https://github.com/user-attachments/assets/82865f7b-b631-47d8-b17b-30ee90eff946)**

## WEB

WEB-интерфейс доступен при запуске приложения(примеры):
1. http://localhost:8080/web/calculate
2. http://localhost:8080/web/expressions
3. http://localhost:8080/web/expressions/{id}
