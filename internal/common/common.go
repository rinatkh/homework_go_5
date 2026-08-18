package common

import "strings"

// TODO: ReadFileOrDefault должен прочитать файл. Если файла нет — вернуть defaultValue без ошибки. Другие ошибки вернуть наружу.
func ReadFileOrDefault(path string, defaultValue string) (string, error) { return "", nil }

// TODO: WriteFileLines должен записать строки в файл, каждую с переводом строки в конце.
func WriteFileLines(path string, lines []string) error { return nil }

// TODO: CountFileLines должен открыть файл, посчитать строки через Scanner и вернуть ошибку открытия или сканирования.
func CountFileLines(path string) (int, error) { return 0, nil }

func Example() string {
	var out strings.Builder
	out.WriteString("common file tasks")
	return out.String()
}
