package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"tunnerse/config"
	"tunnerse/models"

	"time"
	"tunnerse/validation"
)

type TunnelService struct {
	// tunnelRepo *repository.TunnelRepository
	validator *validation.TunnelValidator
	tunnels   map[string]*Tunnel
	mux       sync.Mutex
}

func NewTunnelService() *TunnelService {
	return &TunnelService{
		// tunnelRepo: repository.NewTunnelRepository(),
		validator: validation.NewTunnelValidator(),
		tunnels:   make(map[string]*Tunnel),
	}
}

type Tunnel struct {
	requestCh  chan *http.Request
	responseCh chan []byte
	writerCh   chan http.ResponseWriter
	resetTimer func()
}

func (s *TunnelService) Register(name string) (string, error) {
	err := s.validator.ValidateTunnelRegister(name)
	if err != nil {
		return "", err
	}

	tunnel := models.Tunnel{
		Name:      name,
		CreatedAt: time.Now(),
	}

	// exists := true
	// for exists {
	// 	random := utils.RandomCode(3)
	// 	tunnel.Name = name + "-" + random

	// 	tunnelModel, err := s.tunnelRepo.GetTunnelByName(tunnel.Name)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	if tunnelModel == nil {
	// 		exists = false
	// 	}
	// }

	// err = s.tunnelRepo.Register(&tunnel)
	// if err != nil {
	// 	return "", err
	// }

	t := &Tunnel{
		requestCh:  make(chan *http.Request),
		responseCh: make(chan []byte),
		writerCh:   make(chan http.ResponseWriter),
	}

	inactivityTimer := time.NewTimer(time.Duration(config.AppConfig.TUNNEL_INACTIVITY_LIFE_TIME) * time.Second)
	maxLifetimeTimer := time.NewTimer(time.Duration(config.AppConfig.TUNNEL_LIFE_TIME) * time.Second)

	t.resetTimer = func() {
		if !inactivityTimer.Stop() {
			select {
			case <-inactivityTimer.C:
			default:
			}
		}
		inactivityTimer.Reset(time.Duration(config.AppConfig.TUNNEL_INACTIVITY_LIFE_TIME) * time.Second)
	}

	s.mux.Lock()
	s.tunnels[tunnel.Name] = t
	s.mux.Unlock()

	go func(tunnelName string, t *Tunnel) {
		defer func() {
			s.mux.Lock()
			delete(s.tunnels, tunnelName)
			s.mux.Unlock()
			close(t.requestCh)
			close(t.responseCh)
			close(t.writerCh)
		}()

		select {
		case <-inactivityTimer.C:
		case <-maxLifetimeTimer.C:
		}

	}(tunnel.Name, t)

	return tunnel.Name, nil
}

func (s *TunnelService) Get(name string) ([]byte, error) {
	s.mux.Lock()
	tunnel, exists := s.tunnels[name]
	s.mux.Unlock()

	if !exists {
		return nil, fmt.Errorf("tunnel not found")
	}

	if tunnel.resetTimer != nil {
		tunnel.resetTimer()
	}

	req := <-tunnel.requestCh
	if req == nil {
		return nil, fmt.Errorf("tunnel not found")
	}

	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body.Close()
	}

	sreq := models.SerializableRequest{
		Method: req.Method,
		URL:    req.URL.String(),
		Header: req.Header,
		Body:   string(bodyBytes),
	}

	return json.Marshal(sreq)
}

func (s *TunnelService) Response(name string, body io.ReadCloser) error {
	s.mux.Lock()
	tunnel, exists := s.tunnels[name]
	s.mux.Unlock()

	if !exists {
		return fmt.Errorf("tunnel not found")
	}

	defer body.Close()

	// 🔹 Decodifica o JSON vindo do cliente
	var resp models.ResponseData
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return fmt.Errorf("failed to decode response JSON: %w", err)
	}

	// 🔹 Decodifica o corpo de base64
	bodyDecoded, err := base64.StdEncoding.DecodeString(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode base64 body: %w", err)
	}

	// 🔹 Aguarda o writer
	wr := <-tunnel.writerCh
	if wr == nil {
		return fmt.Errorf("tunnel not found")
	}

	// 🔹 Copia os headers recebidos
	for key, values := range resp.Headers {
		for _, value := range values {
			wr.Header().Add(key, value)
		}
	}

	// 🔹 Escreve o status e o corpo
	wr.WriteHeader(resp.StatusCode)
	_, err = wr.Write(bodyDecoded)
	return err
}

func (s *TunnelService) Tunnel(name string, w http.ResponseWriter, r *http.Request) error {
	err := s.validator.ValidateTunnelRegister(name)
	if err != nil {
		return err
	}

	s.mux.Lock()
	tunnel, exists := s.tunnels[name]
	s.mux.Unlock()

	if !exists {
		return fmt.Errorf("tunnel not found")
	}

	if exists && tunnel.resetTimer != nil {
		tunnel.resetTimer()
	}

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}
	}
	r.Body.Close()

	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	clonedRequest := r.Clone(r.Context())
	clonedRequest.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	select {
	case tunnel.requestCh <- clonedRequest:
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout")
	}

	select {
	case tunnel.writerCh <- w:
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout")
	}

	return nil
}

func (s *TunnelService) Close(name string) error {
	s.mux.Lock()
	tunnel, exists := s.tunnels[name]
	if exists {
		delete(s.tunnels, name)
	}
	s.mux.Unlock()

	if !exists {
		return fmt.Errorf("tunnel not found")
	}

	close(tunnel.requestCh)
	close(tunnel.responseCh)
	close(tunnel.writerCh)

	return nil
}

func (s *TunnelService) NotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)

	path := filepath.Join("static/not_found", "index.html")

	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "404 - tunnel not found", http.StatusNotFound)
		return
	}

	w.Write(data)
}

func (s *TunnelService) Timeout(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)

	path := filepath.Join("static/timeout", "index.html")

	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "404 - tunnel timeout", http.StatusNotFound)
		return
	}

	w.Write(data)
}
