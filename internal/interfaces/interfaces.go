package interfaces

import (
	"fmt"
	"strings"
)

type User struct {
	ID   int
	Name string
}

// TODO: String должен вернуть пользователя в формате "user#<ID>:<Name>".
func (u User) String() string {
	return ""
}

// TODO: FormatOne должен вызвать String у переданного значения. Если значение nil, верни пустую строку.
func FormatOne(value fmt.Stringer) string {
	return ""
}

// TODO: FormatMany должен вернуть строки для всех элементов в исходном порядке. nil элементы пропускаются.
func FormatMany(values []fmt.Stringer) []string {
	return nil
}

type EmailNotifier struct {
	Address string
}

// TODO: Notify должен вернуть "email:<Address>:<message>".
func (n EmailNotifier) Notify(message string) string {
	return ""
}

type PushNotifier struct {
	DeviceID string
}

// TODO: Notify должен вернуть "push:<DeviceID>:<message>".
func (n PushNotifier) Notify(message string) string {
	return ""
}

type Notifier interface {
	Notify(message string) string
}

// TODO: SendNotification должен отправить одно сообщение через Notifier. Если Notifier nil, верни пустую строку.
func SendNotification(n Notifier, message string) string {
	return ""
}

// TODO: SendBatch должен отправить все сообщения и вернуть результаты в том же порядке. nil Notifier даёт пустой результат.
func SendBatch(n Notifier, messages []string) []string {
	return nil
}

type StaticLoader struct {
	Value string
	Err   error
}

// TODO: Load должен вернуть поля Value и Err без дополнительной обработки.
func (l StaticLoader) Load() (string, error) {
	return "", nil
}

type Loader interface {
	Load() (string, error)
}

// TODO: LoadUpper должен загрузить строку, привести её к верхнему регистру и вернуть ошибку без потери, если загрузка не удалась.
func LoadUpper(l Loader) (string, error) {
	return "", nil
}

// TODO: Describe должен вернуть описание динамического типа: nil, string:<value>, user:<String>, notifier:<Notify("ping")>, loader:<loaded>, error:<text> или unknown:<T>.
func Describe(value any) string {
	return ""
}

// TODO: OnlyErrors должен вернуть тексты только ненулевых ошибок в исходном порядке.
func OnlyErrors(values []error) []string {
	return nil
}

// TODO: JoinStringers должен отформатировать ненулевые Stringer и объединить их через sep.
func JoinStringers(values []fmt.Stringer, sep string) string {
	return ""
}

// TODO: FirstNonEmptyStringer должен вернуть первую непустую строку среди Stringer. nil и пустые строки пропускаются.
func FirstNonEmptyStringer(values []fmt.Stringer) string {
	return ""
}

// TODO: CountStringers должен посчитать, сколько элементов реально реализуют fmt.Stringer.
func CountStringers(values []any) int {
	return 0
}

// TODO: NormalizeAndFormat должен обрезать пробелы вокруг String() и схлопнуть пустой результат в "empty".
func NormalizeAndFormat(value fmt.Stringer) string {
	return ""
}

func Example() string {
	var out strings.Builder
	fmt.Fprintf(&out, "user: %s\n", FormatOne(User{ID: 7, Name: "Maria"}))
	fmt.Fprintf(&out, "email: %s\n", SendNotification(EmailNotifier{Address: "team@example.com"}, "hello"))
	fmt.Fprintf(&out, "push: %s\n", SendNotification(PushNotifier{DeviceID: "ios-1"}, "hello"))
	fmt.Fprintf(&out, "describe: %s", Describe(User{ID: 8, Name: "Alex"}))
	return out.String()
}
