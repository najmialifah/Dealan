package service

import (
	"context"
	"log"
	"sync"
	"time"

	"chat-service/models"
	"chat-service/repository"
	"github.com/gorilla/websocket"
)

// Client merepresentasikan koneksi websocket aktif dari seorang pengguna (user/driver).
type Client struct {
	UserID   string
	Role     string // user | driver
	OrderID  string
	Conn     *websocket.Conn
	Send     chan models.WSMessage
	Hub      ChatService
}

// ReadPump membaca pesan masuk dari WebSocket secara terus menerus (looping).
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	// Konfigurasi limit baca websocket untuk keamanan
	c.Conn.SetReadLimit(4096)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg models.WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket-Error] error read websocket: %v", err)
			}
			break
		}

		// Pastikan field diisi dengan benar dari context client
		msg.OrderID = c.OrderID
		msg.SenderID = c.UserID
		msg.SenderRole = c.Role
		msg.SentAt = time.Now()

		// Kirim ke Hub untuk disebarkan (broadcast)
		c.Hub.Broadcast(msg)

		// Simpan pesan ke PostgreSQL secara asinkron agar tidak memblokir websocket
		go func(m models.WSMessage) {
			dbMsg := &models.ChatMessage{
				OrderID:    m.OrderID,
				SenderID:   m.SenderID,
				SenderRole: m.SenderRole,
				Message:    m.Message,
				Type:       m.Type,
				SentAt:     m.SentAt,
				ReadStatus: false,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// Panggil repository langsung melalui helper di service (atau type assert)
			// Kita simpan lewat repository yang ada di concrete service.
			// Karena repository di dalam chatServiceImpl bersifat private, kita bisa
			// menambahkan method SaveMessage di interface ChatService jika dibutuhkan,
			// namun cara paling praktis adalah membuat method SaveMessage di interface.
			if err := c.Hub.SaveMessage(ctx, dbMsg); err != nil {
				log.Printf("[ChatService-Warning] Gagal menyimpan pesan ke database: %v", err)
			}
		}(msg)
	}
}

// WritePump mengirimkan pesan keluar ke WebSocket client secara asinkron.
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(msg); err != nil {
				log.Printf("[WebSocket-Error] gagal mengirim json ke client: %v", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ChatService mendefinisikan interface bisnis untuk chat-service.
type ChatService interface {
	Run()
	Register(c *Client)
	Unregister(c *Client)
	Broadcast(msg models.WSMessage)
	SaveMessage(ctx context.Context, msg *models.ChatMessage) error
	GetHistory(ctx context.Context, orderID string) ([]models.ChatMessage, error)
	CreateRoom(ctx context.Context, orderID, userID, driverID string) error
	GetClientsInRoom(orderID string) int // untuk unit testing
}

type chatServiceImpl struct {
	repo       repository.ChatRepository
	clients    map[string]map[string]*Client // map[order_id]map[user_id]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan models.WSMessage
	mu         sync.RWMutex
}

// NewChatService membuat instance layanan ChatService terintegrasi dengan WebSocket Hub.
func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatServiceImpl{
		repo:       repo,
		clients:    make(map[string]map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan models.WSMessage),
	}
}

func (h *chatServiceImpl) Register(c *Client) {
	h.register <- c
}

func (h *chatServiceImpl) Unregister(c *Client) {
	h.unregister <- c
}

func (h *chatServiceImpl) Broadcast(msg models.WSMessage) {
	h.broadcast <- msg
}

func (h *chatServiceImpl) SaveMessage(ctx context.Context, msg *models.ChatMessage) error {
	return h.repo.SaveMessage(ctx, msg)
}


func (h *chatServiceImpl) GetClientsInRoom(orderID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if room, exists := h.clients[orderID]; exists {
		return len(room)
	}
	return 0
}

// Run menjalankan loop background WebSocket Hub untuk mengelola pendaftaran client dan broadcast pesan.
func (h *chatServiceImpl) Run() {
	log.Println("WebSocket Hub sedang berjalan...")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, exists := h.clients[client.OrderID]; !exists {
				h.clients[client.OrderID] = make(map[string]*Client)
			}
			h.clients[client.OrderID][client.UserID] = client
			h.mu.Unlock()
			log.Printf("[WebSocket] Client terhubung: UserID=%s, Role=%s, OrderID=%s", client.UserID, client.Role, client.OrderID)

		case client := <-h.unregister:
			h.mu.Lock()
			if room, exists := h.clients[client.OrderID]; exists {
				if _, ok := room[client.UserID]; ok {
					delete(room, client.UserID)
					close(client.Send)
					log.Printf("[WebSocket] Client terputus: UserID=%s, OrderID=%s", client.UserID, client.OrderID)
				}
				if len(room) == 0 {
					delete(h.clients, client.OrderID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			roomClients, exists := h.clients[msg.OrderID]
			if exists {
				// Kirim pesan ke semua client di dalam room pesanan yang sama
				for _, client := range roomClients {
					// Kirim ke channel Send masing-masing client
					select {
					case client.Send <- msg:
					default:
						// Jika channel terhambat, bersihkan client secara paksa
						log.Printf("[WebSocket] Antrean penuh untuk UserID=%s, membersihkan client...", client.UserID)
						go func(c *Client) {
							h.unregister <- c
							c.Conn.Close()
						}(client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetHistory mengambil riwayat pesan chat yang tersimpan di database.
func (h *chatServiceImpl) GetHistory(ctx context.Context, orderID string) ([]models.ChatMessage, error) {
	return h.repo.GetChatHistory(ctx, orderID)
}

// CreateRoom membuat room obrolan baru di database PostgreSQL.
func (h *chatServiceImpl) CreateRoom(ctx context.Context, orderID, userID, driverID string) error {
	room := &models.ChatRoom{
		OrderID:  orderID,
		UserID:   userID,
		DriverID: driverID,
		Status:   "active",
	}
	return h.repo.CreateRoom(ctx, room)
}