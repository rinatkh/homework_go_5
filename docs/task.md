# Задания к домашней работе 5

Перед выполнением задач прочитай этот файл и TODO рядом с каждой функцией. Тесты проверяют требования, но не являются единственным источником условия.

## Общие правила

- Не меняй сигнатуры функций и методов.
- Для функций, принимающих `io.Reader`/`io.Writer`, не привязывай решение к файлам.
- Для файловых операций закрывай файлы после успешного открытия.
- Там, где нужно распознать отсутствие файла, используй сравнение ошибки по смыслу, а не строку ошибки.
- `nil` и пустые коллекции обрабатывай так, как указано в таблице конкретной задачи.

## 01 · Интерфейсы

### 1. User.String

**Сигнатура:** `func (u User) String() string`

| Вход | Вывод и правила |
|---|---|
| `User{ID, Name}` | Строка `user#<ID>:<Name>`. Имя не триммится, отрицательный ID сохраняется. |

### 2. FormatOne

**Сигнатура:** `func FormatOne(value fmt.Stringer) string`

| Вход | Вывод и правила |
|---|---|
| Одно значение, реализующее `fmt.Stringer`, или `nil`. | Результат `String()`. Для `nil` — пустая строка. |

### 3. FormatMany

**Сигнатура:** `func FormatMany(values []fmt.Stringer) []string`

| Вход | Вывод и правила |
|---|---|
| Слайс `fmt.Stringer`, внутри могут быть `nil`. | Строки ненулевых элементов в исходном порядке. Для `nil`/пустого входа — пустой слайс. |

### 4. EmailNotifier.Notify

**Сигнатура:** `func (n EmailNotifier) Notify(message string) string`

| Вход | Вывод и правила |
|---|---|
| Адрес получателя и сообщение. | Строка `email:<Address>:<message>` без дополнительной валидации. |

### 5. PushNotifier.Notify

**Сигнатура:** `func (n PushNotifier) Notify(message string) string`

| Вход | Вывод и правила |
|---|---|
| ID устройства и сообщение. | Строка `push:<DeviceID>:<message>` без дополнительной валидации. |

### 6. SendNotification

**Сигнатура:** `func SendNotification(n Notifier, message string) string`

| Вход | Вывод и правила |
|---|---|
| Notifier и сообщение. | Результат `Notify`. Для `nil` Notifier — пустая строка. |

### 7. SendBatch

**Сигнатура:** `func SendBatch(n Notifier, messages []string) []string`

| Вход | Вывод и правила |
|---|---|
| Notifier и список сообщений. | Результаты отправки в исходном порядке. Для `nil` Notifier — пустой слайс. |

### 8. StaticLoader.Load

**Сигнатура:** `func (l StaticLoader) Load() (string, error)`

| Вход | Вывод и правила |
|---|---|
| Поля `Value` и `Err` внутри структуры. | Возвращаются ровно эти поля без преобразования. |

### 9. LoadUpper

**Сигнатура:** `func LoadUpper(l Loader) (string, error)`

| Вход | Вывод и правила |
|---|---|
| Loader, который возвращает строку или ошибку. | При успехе строка в верхнем регистре. При ошибке — пустая строка и исходная ошибка. `nil` loader — ошибка. |

### 10. Describe

**Сигнатура:** `func Describe(value any) string`

| Вход | Вывод и правила |
|---|---|
| Любое значение. | Описание динамического типа: `nil`, `string:<value>`, `user:<String>`, `notifier:<Notify ping>`, `loader:<value>`, `loader-error:<err>`, `error:<err>`, `unknown:<T>`. |

### 11. OnlyErrors

**Сигнатура:** `func OnlyErrors(values []error) []string`

| Вход | Вывод и правила |
|---|---|
| Слайс ошибок, внутри могут быть `nil`. | Тексты ненулевых ошибок в исходном порядке. Для пустого входа — пустой слайс. |

### 12. JoinStringers

**Сигнатура:** `func JoinStringers(values []fmt.Stringer, sep string) string`

| Вход | Вывод и правила |
|---|---|
| Слайс `fmt.Stringer` и разделитель. | Строки ненулевых элементов, объединённые через `sep`. |

### 13. FirstNonEmptyStringer

**Сигнатура:** `func FirstNonEmptyStringer(values []fmt.Stringer) string`

| Вход | Вывод и правила |
|---|---|
| Слайс `fmt.Stringer`. | Первая строка, у которой `String()` не пустая. `nil` и пустые строки пропускаются. |

### 14. CountStringers

**Сигнатура:** `func CountStringers(values []any) int`

| Вход | Вывод и правила |
|---|---|
| Слайс любых значений. | Количество элементов, которые реально реализуют `fmt.Stringer`. |

### 15. NormalizeAndFormat

**Сигнатура:** `func NormalizeAndFormat(value fmt.Stringer) string`

| Вход | Вывод и правила |
|---|---|
| Одно `fmt.Stringer` или `nil`. | Результат `String()` с обрезанными пробелами по краям. Если после этого пусто или вход `nil` — `empty`. |

## 02 · Method set T и *T

### 1. NewCounter

**Сигнатура:** `func NewCounter(name string, value int) Counter`

| Вход | Вывод и правила |
|---|---|
| Имя и начальное значение. | Новый `Counter` с этими полями. |

### 2. Counter.Snapshot

**Сигнатура:** `func (c Counter) Snapshot() int`

| Вход | Вывод и правила |
|---|---|
| Counter. | Текущее значение `Value`; объект не меняется. |

### 3. Counter.Label

**Сигнатура:** `func (c Counter) Label() string`

| Вход | Вывод и правила |
|---|---|
| Counter. | Строка `<Name>=<Value>`. |

### 4. (*Counter).Add

**Сигнатура:** `func (c *Counter) Add(delta int)`

| Вход | Вывод и правила |
|---|---|
| Указатель на Counter и delta. | Исходный объект увеличивается на delta. `nil` receiver не паникует. |

### 5. (*Counter).Reset

**Сигнатура:** `func (c *Counter) Reset()`

| Вход | Вывод и правила |
|---|---|
| Указатель на Counter. | Исходный объект получает `Value=0`. `nil` receiver не паникует. |

### 6. UseSnapshoter

**Сигнатура:** `func UseSnapshoter(s Snapshoter) int`

| Вход | Вывод и правила |
|---|---|
| Интерфейс с методом `Snapshot()`. | Snapshot или `0`, если интерфейс nil. |

### 7. UseLabeler

**Сигнатура:** `func UseLabeler(l Labeler) string`

| Вход | Вывод и правила |
|---|---|
| Интерфейс с методом `Label()`. | Label или пустая строка, если интерфейс nil. |

### 8. UseAdder

**Сигнатура:** `func UseAdder(a Adder, delta int) int`

| Вход | Вывод и правила |
|---|---|
| Интерфейс с `Add(delta int)`. | Вызывает `Add`. Если объект также умеет `Snapshot`, возвращает новое значение; иначе `0`. |

### 9. IsSnapshoter

**Сигнатура:** `func IsSnapshoter(value any) bool`

| Вход | Вывод и правила |
|---|---|
| Любое значение. | true, если динамический тип реализует `Snapshoter`. |

### 10. IsAdder

**Сигнатура:** `func IsAdder(value any) bool`

| Вход | Вывод и правила |
|---|---|
| Любое значение. | true, если динамический тип реализует `Adder`. |

### 11. CloneAndAdd

**Сигнатура:** `func CloneAndAdd(c Counter, delta int) Counter`

| Вход | Вывод и правила |
|---|---|
| Counter и delta. | Новая копия с увеличенным значением; исходный Counter не меняется. |

### 12. AddInPlace

**Сигнатура:** `func AddInPlace(c *Counter, delta int) int`

| Вход | Вывод и правила |
|---|---|
| Указатель на Counter и delta. | Меняет исходный объект и возвращает новое значение. `nil` -> `0`. |

### 13. Profile.Display

**Сигнатура:** `func (p Profile) Display() string`

| Вход | Вывод и правила |
|---|---|
| Profile. | Строка `<Name>(<Age>)`. |

### 14. (*Profile).Rename

**Сигнатура:** `func (p *Profile) Rename(name string)`

| Вход | Вывод и правила |
|---|---|
| Указатель на Profile и новое имя. | Меняет `Name` у исходного объекта. `nil` receiver не паникует. |

### 15. IsRenamer

**Сигнатура:** `func IsRenamer(value any) bool`

| Вход | Вывод и правила |
|---|---|
| Любое значение. | true, если динамический тип реализует `Renamer`. |

## 03 · io.Reader / io.Writer

### 1. ReadAllText

**Сигнатура:** `func ReadAllText(r io.Reader) (string, error)`

| Вход | Вывод и правила |
|---|---|
| Любой `io.Reader`. | Весь поток как строка или ошибка чтения. |

### 2. CountBytes

**Сигнатура:** `func CountBytes(r io.Reader) (int, error)`

| Вход | Вывод и правила |
|---|---|
| Любой `io.Reader`. | Количество байт в потоке. Для ошибки чтения вернуть ошибку. |

### 3. WriteString

**Сигнатура:** `func WriteString(w io.Writer, text string) (int, error)`

| Вход | Вывод и правила |
|---|---|
| Writer и строка. | Записать строку и вернуть количество байт плюс ошибку записи. |

### 4. WriteLines

**Сигнатура:** `func WriteLines(w io.Writer, lines []string) error`

| Вход | Вывод и правила |
|---|---|
| Writer и строки. | Каждую строку записать с `\n`. При первой ошибке остановиться. |

### 5. CopyAll

**Сигнатура:** `func CopyAll(dst io.Writer, src io.Reader) (int64, error)`

| Вход | Вывод и правила |
|---|---|
| Источник и приёмник байт. | Скопировать весь поток и вернуть число байт. |

### 6. CopyUpper

**Сигнатура:** `func CopyUpper(dst io.Writer, src io.Reader) (int64, error)`

| Вход | Вывод и правила |
|---|---|
| Reader и Writer. | Прочитать текст, привести к верхнему регистру и записать в Writer. |

### 7. ScanLines

**Сигнатура:** `func ScanLines(r io.Reader) ([]string, error)`

| Вход | Вывод и правила |
|---|---|
| Reader с текстом. | Все строки без символов перевода строки. После цикла учитывать ошибку Scanner. |

### 8. CountNonEmptyLines

**Сигнатура:** `func CountNonEmptyLines(r io.Reader) (int, error)`

| Вход | Вывод и правила |
|---|---|
| Reader с текстом. | Количество строк, которые после `TrimSpace` не пустые. |

### 9. FirstLine

**Сигнатура:** `func FirstLine(r io.Reader) (string, error)`

| Вход | Вывод и правила |
|---|---|
| Reader с текстом. | Первая строка. Если строк нет — пустая строка без ошибки. |

### 10. ReadCSVLike

**Сигнатура:** `func ReadCSVLike(r io.Reader) ([][]string, error)`

| Вход | Вывод и правила |
|---|---|
| Текстовые строки с ячейками через запятую. | Таблица строк; у каждой ячейки убраны пробелы по краям. |

### 11. LimitRead

**Сигнатура:** `func LimitRead(r io.Reader, limit int) (string, error)`

| Вход | Вывод и правила |
|---|---|
| Reader и лимит байт. | Не больше `limit` байт. Отрицательный лимит считать нулём. |

### 12. RepeatToWriter

**Сигнатура:** `func RepeatToWriter(w io.Writer, text string, count int) error`

| Вход | Вывод и правила |
|---|---|
| Writer, строка и количество повторов. | Записать строку `count` раз подряд. `count <= 0` ничего не пишет. |

### 13. BufferReport

**Сигнатура:** `func BufferReport(lines []string) string`

| Вход | Вывод и правила |
|---|---|
| Слайс строк. | Строки формата `<index>:<line>\n`, собранные через Writer/Buffer. |

### 14. ReadAndTrim

**Сигнатура:** `func ReadAndTrim(r io.Reader) (string, error)`

| Вход | Вывод и правила |
|---|---|
| Reader. | Весь текст без пробелов по краям или ошибка чтения. |

### 15. WriteKeyValues

**Сигнатура:** `func WriteKeyValues(w io.Writer, values map[string]string) error`

| Вход | Вывод и правила |
|---|---|
| Writer и map строк. | Пары в отсортированном по ключу порядке, формат `key=value\n`. |

## 04 · Общие файловые задачи

### 1. ReadFileOrDefault

**Сигнатура:** `func ReadFileOrDefault(path string, defaultValue string) (string, error)`

| Вход | Вывод и правила |
|---|---|
| Путь к файлу и значение по умолчанию. | Если файл есть — содержимое. Если файла нет — default без ошибки. Другие ошибки вернуть. |

### 2. WriteFileLines

**Сигнатура:** `func WriteFileLines(path string, lines []string) error`

| Вход | Вывод и правила |
|---|---|
| Путь и строки. | Создать/перезаписать файл, каждую строку записать с `\n`, закрыть файл. |

### 3. CountFileLines

**Сигнатура:** `func CountFileLines(path string) (int, error)`

| Вход | Вывод и правила |
|---|---|
| Путь к файлу. | Открыть файл, посчитать строки через Scanner, закрыть файл, вернуть ошибки открытия/сканирования. |

---

# Самопроверка перед отправкой

```bash
make compile
make fmt
make fmt-check
make test-interfaces
make test-methodsets
make test-io
make test-common
make test-unit
make test-integration
make ci
```

Перед PR проверь:

- все TODO реализованы;
- сигнатуры и тесты не изменены;
- каждый раздел проходит отдельно;
- `make ci` проходит полностью;
- в репозиторий не добавлены бинарники, coverage-файлы и файлы IDE.
