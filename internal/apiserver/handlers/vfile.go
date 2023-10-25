package handlers

import (
	"fmt"
	"go-gost/internal/apiserver/response"
	"go-gost/internal/lib/signature"
	"io"
	"net/http"
)

// Проверка ЗЦП с загрузкой файлов
func (h *Handler) VerifyFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Error(err.Error())
	}

	documetFile, _, err := r.FormFile("document")
	if err != nil {
		h.logger.Error("Ошибка загрузки файла документа.", err.Error())
		err = fmt.Errorf("Ошибка загрузки файла документа.")
		response.ErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer documetFile.Close()

	documentContent, err := io.ReadAll(documetFile)
	if err != nil {
		h.logger.Error("Невозможно прочитать содержимое документа.", err.Error())
		err = fmt.Errorf("Невозможно прочитать содержимое документа.")
		response.ErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	signFile, _, err := r.FormFile("signature")
	if err != nil {
		h.logger.Error("Ошибка загрузки файла подписи.", err.Error())
		err = fmt.Errorf("Ошибка загрузки файла подписи.")
		response.ErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	defer signFile.Close()

	signContent, err := io.ReadAll(signFile)
	if err != nil {
		h.logger.Error("Невозможно прочитать содержимое файла подписи.", err.Error())
		err = fmt.Errorf("Невозможно прочитать содержимое файла подписи.")
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
