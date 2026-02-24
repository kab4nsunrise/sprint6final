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


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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


func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	/
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Ошибка при парсинге формы: %v", err)
		http.Error(w, "Ошибка при обработке формы", http.StatusInternalServerError)
		return
	}

	
	if r.MultipartForm == nil {
		log.Printf("MultipartForm is nil")
		http.Error(w, "Ошибка при обработке формы", http.StatusInternalServerError)
		return
	}

	
	var keys []string
	for k := range r.MultipartForm.File {
		keys = append(keys, k)
	}
	log.Printf("Ключи в MultipartForm.File: %v", keys)

	
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		log.Printf("Нет файла с именем 'file'")
		http.Error(w, "Ошибка при получении файла", http.StatusInternalServerError)
		return
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		http.Error(w, "Ошибка при открытии файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	
	converted, err := service.DetectAndConvert(string(content))
	if err != nil {
		log.Printf("Ошибка при конвертации: %v", err)
		http.Error(w, "Ошибка при конвертации файла", http.StatusInternalServerError)
		return
	}

	
	timestamp := time.Now().UTC().Format("20060102-150405")
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		if service.IsMorseCode(string(content)) {
			ext = ".txt"
		} else {
			ext = ".morse"
		}
	}
	newFilename := "converted_" + timestamp + ext

	
	newFile, err := os.Create(newFilename)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		http.Error(w, "Ошибка при создании файла результата", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = newFile.WriteString(converted)
	if err != nil {
		log.Printf("Ошибка при записи в файл: %v", err)
		http.Error(w, "Ошибка при записи результата", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(converted))
}
