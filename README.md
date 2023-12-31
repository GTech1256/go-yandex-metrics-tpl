# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

### Локальный запуск автотестов

1. Скомпилируйте ваши сервер и агент в папках `cmd/server` и `cmd/agent` командами `go build -o server *.go` и `go build -o agent *.go` соответственно.
2. Скачайте [бинарный файл с автотестами](https://github.com/Yandex-Practicum/go-autotests/releases/latest) для вашей ОС — например, `metricstest-darwin-arm64` для MacOS на процессоре Apple Silicon.
3. Разместите бинарный файл так, чтобы он был доступен для запуска из командной строки, — пропишите путь в переменную `$PATH`.
4. Ознакомьтесь с параметрами запуска автотестов в файле `.github/workflows/metricstest.yml` вашего репозитория. Автотесты для разных инкрементов требуют различных аргументов для запуска.

Пример запуска теста для первого инкремента:

```shell
metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent
```

### Запуск на Mac с процессором Apple Silicon

Если у вас возникают трудности с локальным запуском автотестов на компьютере Mac на базе процессора Apple Silicon (M1 и старше), убедитесь, что:

- Вы запускаете бинарный файл с суффиксом `darwin-arm64`.
- У файла выставлен флаг на исполнение командой `chmod +x <filename>`.
- Вы разрешили выполнение неподписанного бинарного файл в разделе «Безопасность» настроек системы.

### Разрешение запуска неподписанного кода на Mac с процессором Apple Silicon

1. Зайдите в раздел «Безопасность» настроек системы и найдите раздел, в котором написано: «Использование <имя файла> было заблокировано, так как он не от идентифицированного разработчика». Этот раздел появляется только после попытки запустить неподписанный бинарный файл.
2. Нажмите на кнопку «Разрешить» (Allow anyway), чтобы внести бинарный файл в список разрешённых к запуску.

<img width="1440" alt="30" src="https://user-images.githubusercontent.com/85521342/228195019-89767be7-a7e5-4f07-b867-baf2ce8344e8.png">

3. Введите пароль или приложите палец к Touch ID для подтверждения разрешения.

<img width="1440" alt="31" src="https://user-images.githubusercontent.com/85521342/228199358-f9e0dbf7-e7ea-4be8-b2f4-e39f6f1bdfc2.png">

Операция выполняется однократно для каждого бинарного файла, который вы хотите запустить. Данную политику безопасности невозможно отключить на компьютерах Mac с процессором Apple Silicon.

#### Запуск

```shell
./metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=/Users/ribakakin/Desktop/web/go-practicum/go-yandex-metrics-tpl/cmd/memory/memory
```
```shell
./autotests/metricstest-darwin-arm64 -test.v -test.run=^TestIteration2[AB]*$ \
  -source-path=/Users/ribakakin/Desktop/web/go-practicum/go-yandex-metrics-tpl \
  -agent-binary-path=/Users/ribakakin/Desktop/web/go-practicum/go-yandex-metrics-tpl/cmd/agent/agent
```

### Убить процесс, слушающий определенный порт

`lsof -i:8080`
`kill -9 <pid>`

```shell
kill -9 $(lsof -t -i:8080)
```

### Сборка
```shell
go build -o cmd/agent cmd/agent/main.go
mv cmd/agent/main cmd/agent/agent
```

```shell
go build -o cmd/memory cmd/memory/main.go
mv cmd/memory/main cmd/memory/memory
```

### Миграции

Добавление миграции
`migrate create -ext sql -dir internal/server/config/db/migrations -seq create_users_table`


