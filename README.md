# Telegram Bot с Emoji Captcha

Этот бот для Telegram использует эмодзи-капчу для проверки новых участников группы.

## Настройка

1. Клонируйте репозиторий
2. Создайте файл `.env` в корневой директории проекта
3. Заполните `.env` файл, используя пример ниже

### Пример .env файла

Токен вашего Telegram бота
TELEGRAM*BOT_TOKEN="ваш*токен_бота"

Количество эмодзи в капче
EMOJI_COUNT=5

Количество кнопок с эмодзи
EMOJI_BUTTON_COUNT=10

Время на прохождение капчи (в минутах)
CAPTCHA_TIMEOUT_MINUTES=5

Приветственное сообщение
%s - имя пользователя, %d - время на капчу, %s - капча
WELCOME_MESSAGE="Добро пожаловать, %s! У вас есть %d минут, чтобы ввести следующую капчу, нажимая на кнопки в правильном порядке:\n%s"

Сообщение при успешном прохождении капчи
CAPTCHA_SUCCESS_MESSAGE="Капча пройдена успешно! Теперь вы можете писать в чате."

Сообщение при частичном успехе (правильный выбор, но капча еще не завершена)
CAPTCHA_PARTIAL_SUCCESS_MESSAGE="Правильно! Продолжайте."

Сообщение при неправильном выборе
CAPTCHA_FAIL_MESSAGE="Неправильный порядок. Попробуйте еще раз."

Объявление о прохождении капчи
%s - имя пользователя
CAPTCHA_PASSED_ANNOUNCEMENT="Пользователь @%s успешно прошел капчу!"

Сообщение при неудачной попытке удаления пользователя
KICK_FAIL_MESSAGE="Я не справился с удалением пользователя, не прошедшего капчу."

Сообщение при успешном удалении пользователя
%d - ID пользователя
KICK_SUCCESS_MESSAGE="Пользователь с ID %d был удален из-за непрохождения капчи."

Пользовательский список эмодзи (необязательно)
Если не указан, будет использован стандартный набор эмодзи
CUSTOM_EMOJI_LIST="🍎,🍐,🍊,🍋,🍌,🍉,🍇,🍓,🫐,🍈,🍒,🍑,🥭,🍍,🥥,🥝,🍅,🍆,🥑"

## Запуск

1. Убедитесь, что у вас установлен Go
2. Выполните команду `go run main.go`

## Примечания

- Убедитесь, что ваш бот имеет права администратора в группе
- Настройте значения в `.env` файле в соответствии с вашими предпочтениями
- Вы можете изменить список эмодзи, отредактировав `CUSTOM_EMOJI_LIST` в `.env` файле

## Пример кода .env

```
TELEGRAM_BOT_TOKEN="Ваш токен"
EMOJI_COUNT=5
EMOJI_BUTTON_COUNT=10
CAPTCHA_TIMEOUT_MINUTES=5
WELCOME_MESSAGE="Добро пожаловать, %s! У вас есть %d минут, чтобы ввести следующую капчу, нажимая на кнопки в правильном порядке:\n%s"
CAPTCHA_SUCCESS_MESSAGE="Капча пройдена успешно! Теперь вы можете писать в чате."
CAPTCHA_PARTIAL_SUCCESS_MESSAGE="Правильно! Продолжайте."
CAPTCHA_FAIL_MESSAGE="Неправильный порядок. Попробуйте еще раз."
CAPTCHA_PASSED_ANNOUNCEMENT="Пользователь @%s успешно прошел капчу!"
KICK_FAIL_MESSAGE="Я не справился с удалением пользователя, не прошедшего капчу."
KICK_SUCCESS_MESSAGE="Пользователь с ID %d был удален из-за непрохождения капчи."
CUSTOM_EMOJI_LIST="🍎,🍐,🍊,🍋,🍌,🍉,🍇,🍓,🫐,🍈,🍒,🍑,🥭,🍍,🥥,🥝,🍅,🍆,🥑"
```
