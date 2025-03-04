// Package server предоставляет утилиты HTTP-ответов для сервера.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"

	"go.uber.org/zap"
)

// writeJSON записывает ответ в формате JSON с заданным статусом и значением.
func (s *Server) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set(common.HeaderContentType, common.ContentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.logger.Error("Failed to write JSON response", zap.Error(err))
	}
}

// writeError пишет ответ об ошибке с заданным статусом и сообщением.
func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set(common.HeaderContentType, common.ContentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		s.logger.Error("Failed to write error response", zap.Error(err))
	}
}
