package web_service

import (
	"fmt"
)

func (s *WebService) GetMainPage() ([]byte, error) {
	html, err := s.webRepository.GetFile("index.html")
	if err != nil {
		return nil, fmt.Errorf("get file from repository: %w", err)
	}

	return html, nil
}
