package service

import (
	"Museum/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type AiService struct {
	repo   *repository.ExhibitRepository
	client *http.Client
	addr   string
}

func NewAiService(repo *repository.ExhibitRepository, client *http.Client, addr string) *AiService {
	return &AiService{repo: repo, client: client, addr: addr}
}

func (s *AiService) AiResponse(ctx context.Context, requestText string) (string, error) {
	//aiReq := struct {
	//	Text string `json:"text"`
	//}{requestText}
	//aiReqBytes, err := json.Marshal(aiReq)
	//if err != nil {
	//	return "", err
	//}
	//resp, err := s.client.Post(s.addr+"/api/v1/embed/text", "application/json; charset=utf-8", bytes.NewReader(aiReqBytes))
	//if err != nil {
	//	return "", err
	//}
	//if resp.StatusCode != http.StatusOK {
	//	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	//	return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	//}
	//defer resp.Body.Close()
	//aiResp := struct {
	//	Vector []float32 `json:"vector"`
	//}{}
	//if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
	//	return "", fmt.Errorf("decode response: %w", err)
	//}
	//list, err := s.repo.MostSuitable(ctx, aiResp.Vector)
	list, _, err := s.repo.List(ctx, 0, 100)
	if err != nil {
		return "", fmt.Errorf("not found most suitable exhibits: %w", err)
	}
	type exh struct {
		ID           int64     `json:"id"`
		Category     string    `json:"category"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		CreationDate time.Time `json:"creation_date"`
	}
	listExh := make([]exh, 0, len(list))
	for _, e := range list {
		listExh = append(listExh, exh{ID: e.ID, Category: e.Category, Name: e.Name, Description: e.Description, CreationDate: e.CreationDate})
	}
	chatBotReq := struct {
		Query   string `json:"query"`
		Context []exh  `json:"context"`
	}{requestText, listExh}
	chatBotRespBytes, err := json.Marshal(&chatBotReq)
	if err != nil {
		return "", err
	}
	chatBorRespBytes, err := s.client.Post(s.addr+"/api/v1/rag/generate", "application/json; charset=utf-8", bytes.NewReader(chatBotRespBytes))
	if err != nil {
		slog.Debug("rag response", "err", err)
		return "", err
	}

	defer chatBorRespBytes.Body.Close()
	if chatBorRespBytes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(chatBorRespBytes.Body, 4096))
		return "", fmt.Errorf("unexpected status %d: %s", chatBorRespBytes.StatusCode, body)
	}
	chatBotResp := struct {
		Answer string `json:"answer"`
	}{}
	if err := json.NewDecoder(chatBorRespBytes.Body).Decode(&chatBotResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return chatBotResp.Answer, nil
}
