// internal/handlers/handlers.go
package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// HomeHandler обрабатывает корневой маршрут и отображает HTML-форму
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что это GET запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим и выполняем шаблон
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Printf("Ошибка при парсинге шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Ошибка при выполнении шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

// UploadHandler обрабатывает загрузку файла и его конвертацию
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart форму с максимальным размером 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Ошибка при парсинге формы: %v", err)
		http.Error(w, "Ошибка при обработке формы", http.StatusInternalServerError)
		return
	}

	// Получаем файл из формы
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		http.Error(w, "Ошибка при получении файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	// Конвертируем содержимое
	converted, err := service.DetectAndConvert(string(content))
	if err != nil {
		log.Printf("Ошибка при конвертации: %v", err)
		http.Error(w, "Ошибка при конвертации файла", http.StatusInternalServerError)
		return
	}

	// Создаем имя для нового файла
	timestamp := time.Now().UTC().Format("20060102-150405")
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		// Если расширения нет, добавляем .txt для текста или .morse для кода Морзе
		if service.IsMorseCode(string(content)) {
			ext = ".txt"
		} else {
			ext = ".morse"
		}
	}
	newFilename := "converted_" + timestamp + ext

	// Создаем новый файл
	newFile, err := os.Create(newFilename)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		http.Error(w, "Ошибка при создании файла результата", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Записываем конвертированное содержимое
	_, err = newFile.WriteString(converted)
	if err != nil {
		log.Printf("Ошибка при записи в файл: %v", err)
		http.Error(w, "Ошибка при записи результата", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(converted))
}
