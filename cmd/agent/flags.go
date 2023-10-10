package main

import "flag"

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	flag.StringVar(&port, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&pollInterval, "p", 2, "the frequency of polling metrics")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
