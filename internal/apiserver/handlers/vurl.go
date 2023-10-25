package handlers

import (
	"errors"
	"fmt"
	"go-gost/internal/apiserver/response"
	"go-gost/internal/lib/signature"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/render"
)

// Структура запроса проверки документа с удаленного сервера
type request struct {
	Document  string `json:"document" validate:"required,url"`
	Signature string `json:"signature" validate:"required,url"`
	Referer   string `json:"referer,omitempty"`
}

// Проверка ЭЦП с удаленного сервера
func (h *Handler) VerifyURL(w http.ResponseWriter, r *http.Request) {
	var req request
	client := &http.Client{
		Timeout: h.config.Timeout,
	}
	headerContentTtype := r.Header.Get("Content-Type")

	if !strings.Contains(headerContentTtype, "application/json") {
		h.logger.Error("Содерживое запроса имеет формат отличный от application/json")
		response.ErrorResponse(w, "Содерживое запроса имеет формат отличный от application/json", http.StatusUnsupportedMediaType)
		return
	}

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		h.logger.Error("Тело запроса не содержит данных")
		response.ErrorResponse(w, "Тело запроса не содержит данных", http.StatusBadRequest)
		return
	}

	document, err := url.ParseRequestURI(req.Document)
	if err != nil {
		h.logger.Error("Тело запроса не содержит ссылки на документ или ссылка неверна")
		response.ErrorResponse(w, "Тело запроса не содержит ссылки на документ или ссылка неверна", http.StatusBadRequest)
		return
	}

	requestDocumentFile, err := http.NewRequest("GET", document.String(), nil)
	if err != nil {
		h.logger.Error("Невозможно созать HTTP запрос на скачивание документа")
		response.ErrorResponse(w, "Ошибка получения документа", http.StatusUnprocessableEntity)
		return
	}
	requestDocumentFile.Header.Set("User-Agent", h.config.UserAgent)
	if req.Referer != "" {
		requestDocumentFile.Header.Set("Referer", req.Referer)
	}

	documentFile, err := client.Do(requestDocumentFile)
	if err != nil {
		h.logger.Error("Ошибка получения документа", err.Error(), documentFile.StatusCode)
		response.ErrorResponse(w, "Ошибка получения документа", http.StatusUnprocessableEntity)
		return
	}
	defer documentFile.Body.Close()

	documentContent, err := io.ReadAll(documentFile.Body)
	if err != nil {
		h.logger.Error("Ошибка чтения содержимого документа.", err.Error())
		err = fmt.Errorf("Ошибка чтения содержимого документа.")
		response.ErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	sig, err := url.ParseRequestURI(req.Signature)
	if err != nil {
		h.logger.Error("Тело запроса не содержит ссылки на файл подптси или ссылка неверна")
		response.ErrorResponse(w, "Тело запроса не содержит ссылки на файл подптси или ссылка неверна", http.StatusBadRequest)
		return
	}

	requestSingFile, err := http.NewRequest("GET", sig.String(), nil)
	if err != nil {
		h.logger.Error("Невозможно созать HTTP запрос на скачивание ЭЦП")
		response.ErrorResponse(w, "Ошибка получения файла подписи", http.StatusUnprocessableEntity)
		return
	}
	requestSingFile.Header.Set("User-Agent", h.config.UserAgent)
	if req.Referer != "" {
		requestSingFile.Header.Set("Referer", req.Referer)
	}

	signFile, err := client.Do(requestSingFile)
	if err != nil {
		h.logger.Error("Ошибка получения файла подписи", err.Error(), documentFile.StatusCode)
		response.ErrorResponse(w, "Ошибка файла подписи документа", http.StatusUnprocessableEntity)
		return
	}
	defer signFile.Body.Close()

	signContent, err := io.ReadAll(signFile.Body)
	if err != nil {
		h.logger.Error("Ошибка чтения содержимого ЭЦП.", err.Error())
		err = fmt.Errorf("Ошибка чтения содержимого подписи.")
		response.ErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	sign, err := signature.NewSign(documentContent, signContent)
	if err != nil {
		h.logger.Error("Ошибка при создании экземпляра signature.", err.Error())
		err = fmt.Errorf("Неудалось идентифицировать формат подписи.")
		response.ErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	body := sign.GetReport()

	if err = response.DataResponse(w, body, http.StatusOK); err != nil {
		h.logger.Error(err.Error())
		response.ErrorResponse(w, err.Error(), http.StatusBadRequest)
	}
	return
}
