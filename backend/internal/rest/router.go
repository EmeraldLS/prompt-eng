package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/emeraldls/fyp/internal/service"
	"github.com/emeraldls/fyp/internal/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Config struct {
	ListAddr           string
	EntityExtractorURL string
	ClientOrigin       string
}

type Server struct {
	Config
	mut               sync.RWMutex
	subscribedClients map[string]chan []byte
	// clientsMsgs       []chan []byte
	options serverOptions
}

type serverOption func(*serverOptions)

type serverOptions struct {
	Domain string
	Secure bool
}

func WithDomain(domain string) serverOption {
	return func(so *serverOptions) {
		so.Domain = domain
	}
}

func WithSecure(secure bool) serverOption {
	return func(so *serverOptions) {
		so.Secure = secure
	}
}

func NewServer(cfg Config, opts ...serverOption) *Server {

	defaultServerOps := serverOptions{
		Domain: "localhost",
		Secure: false,
	}

	for _, opt := range opts {
		opt(&defaultServerOps)
	}

	return &Server{
		Config:            cfg,
		options:           defaultServerOps,
		subscribedClients: make(map[string]chan []byte),
	}
}

func (s *Server) SetupRouter() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{s.ClientOrigin},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Type"},
	}))

	r.POST("/chat", s.SendMessage)
	r.GET("/messages", s.MessagesHandler)

	r.Run(s.ListAddr)
}

func extractEntities(req types.ChatMessage, url string) (*types.Entity, error) {
	requestBody, _ := json.Marshal(req)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, errors.New(string(body))
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var entities types.Entity
	json.Unmarshal(body, &entities)

	return &entities, nil
}

// TODO: use jwt for token
func (s *Server) SendMessage(c *gin.Context) {

	bearerHeader := c.GetHeader("Authorization")
	bearer := strings.Split(bearerHeader, " ")
	if len(bearer) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": "Bearer token not provided",
		})
		return
	}

	sessionID := bearer[1]

	var req types.ChatMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": "Invalid input data",
		})
		return
	}

	s.mut.RLock()
	_, ok := s.subscribedClients[sessionID]
	s.mut.RUnlock()

	if !ok {
		fmt.Println("SendMessage: client not subscribed:", sessionID)
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": "Client not subscribed",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": err.Error(),
		})
		return
	}

	entities, err := extractEntities(req, s.EntityExtractorURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": "Unable to extract entities: " + err.Error(),
		})
		return
	}

	rs := service.NewRouteService(os.Getenv("HERE_API_KEY"))
	fromResp, err := rs.EncodeText(entities.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": err.Error(),
		})
		return
	}

	toResp, err := rs.EncodeText(entities.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": err.Error(),
		})
		return
	}

	fromCoordinates := fromResp.Items[0].Position
	toCoordinates := toResp.Items[0].Position

	route, err := rs.GetRoute(types.CAR, []float64{fromCoordinates.Lat, fromCoordinates.Lng}, []float64{toCoordinates.Lat, toCoordinates.Lng})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": err.Error(),
		})
		return
	}

	// go func() {
	// 	llm_sync.UseGoogle(req.Prompt, msgChan)
	// }()

	c.JSON(http.StatusOK, gin.H{
		"STATUS":  "SUCCESS",
		"MESSAGE": route,
	})

}

func (s *Server) MessagesHandler(c *gin.Context) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"STATUS":  "FAILURE",
			"MESSAGE": "SSE not supported",
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	msgChan := make(chan []byte)
	sessionID := uuid.New().String()

	s.mut.Lock()
	s.subscribedClients[sessionID] = msgChan
	s.mut.Unlock()

	fmt.Println("MessagesHandler: client subscribed:", sessionID)

	err := sendSSEEvent(c.Writer, "connected", []byte(sessionID), flusher)
	if err != nil {
		log.Printf("Error sending initial message: %v", err)
		return
	}

	for {
		select {
		case msgChunk, ok := <-msgChan:
			if !ok {
				sendSSEEvent(c.Writer, "close", []byte("Channel closed"), flusher)
				return
			}

			if len(msgChunk) > 0 {
				err := sendSSEEvent(c.Writer, "message", msgChunk, flusher)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					return
				}
			}

		case <-c.Request.Context().Done():
			s.mut.Lock()
			delete(s.subscribedClients, sessionID)
			s.mut.Unlock()
			fmt.Println("MessagesHandler: client unsubscribed:", sessionID)
			return
		}
	}

}

func sendSSEEvent(w io.Writer, event string, data []byte, flusher http.Flusher) error {
	resp, err := formatChatSentEvent(event, data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, resp)
	if err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func formatChatSentEvent(event string, chunk []byte) (string, error) {
	eventData := map[string]interface{}{
		"data": string(chunk),
	}

	buffer := bytes.NewBuffer([]byte{})

	err := json.NewEncoder(buffer).Encode(eventData)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buffer.String()))

	return sb.String(), nil
}
