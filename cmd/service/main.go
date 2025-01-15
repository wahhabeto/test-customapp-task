package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wahhabeto/test-customapp-task/internal/di"
	"github.com/wahhabeto/test-customapp-task/pkg/closer"
	"net/http"
)

func main() {
	// Используем флаг для получения значения rtp из командной строки
	rtp := flag.Float64("rtp", 0.5, "The RTP value between 0 and 1")
	flag.Parse()

	// Проверяем, чтобы значение rtp было в пределах от 0 до 1
	if *rtp <= 0 || *rtp > 1.0 {
		log.Fatal("RTP must be between 0 and 1.0")
	}

	// Ожидаем корректное завершение приложения с помощью defer
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	ctx, cancel := context.WithCancel(context.Background())

	// Добавляем функцию отмены в closer для корректного завершения
	closer.Add(func() error {
		cancel()
		return nil
	})

	// Настроим InsecureSkipVerify для TLS
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Создаем DI контейнер и передаем rtp в качестве аргумента
	container := di.MustCreateDiContainer(*rtp)

	// Инжектируем зависимости и запускаем приложение
	err := container.Invoke(func(a *di.App) error {
		return a.Run(ctx)
	})

	// Если возникла ошибка, выводим её в лог
	if err != nil {
		log.Error("Error running app", "error", err)
	}
}
