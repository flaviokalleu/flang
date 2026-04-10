package whatsapp

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Client wraps the whatsmeow client for Flang.
type Client struct {
	mu        sync.Mutex
	client    *whatsmeow.Client
	connected bool
	dbPath    string
}

// Novo creates a new WhatsApp client.
func Novo(dbPath string) *Client {
	if dbPath == "" {
		dbPath = "whatsapp.db"
	}
	return &Client{dbPath: dbPath}
}

// Conectar initializes the connection and shows QR code if needed.
func (c *Client) Conectar() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx := context.Background()

	dbLog := waLog.Stdout("WA-DB", "WARN", true)
	container, err := sqlstore.New(ctx, "sqlite3",
		"file:"+c.dbPath+"?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("erro ao criar store WhatsApp: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter device: %w", err)
	}

	clientLog := waLog.Stdout("WA", "WARN", true)
	c.client = whatsmeow.NewClient(deviceStore, clientLog)

	if c.client.Store.ID == nil {
		// Need to login with QR code
		fmt.Println("[whatsapp] Escaneie o QR Code para conectar:")
		fmt.Println()

		qrChan, _ := c.client.GetQRChannel(ctx)
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("erro ao conectar WhatsApp: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				// Print QR code as text for terminal
				printQR(evt.Code)
			} else if evt.Event == "success" {
				fmt.Println("[whatsapp] ✓ Conectado com sucesso!")
				break
			} else if evt.Event == "timeout" {
				return fmt.Errorf("timeout ao esperar QR code")
			}
		}
	} else {
		// Already logged in
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("erro ao conectar WhatsApp: %w", err)
		}
		fmt.Println("[whatsapp] ✓ Reconectado")
	}

	c.connected = true
	return nil
}

// EnviarMensagem sends a text message to a phone number.
func (c *Client) EnviarMensagem(telefone string, mensagem string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected || c.client == nil {
		return fmt.Errorf("WhatsApp não conectado")
	}

	// Clean phone number
	phone := limparTelefone(telefone)
	if phone == "" {
		return fmt.Errorf("telefone inválido: %s", telefone)
	}

	jid := types.NewJID(phone, types.DefaultUserServer)

	msg := &waE2E.Message{
		Conversation: proto.String(mensagem),
	}

	ctx := context.Background()
	_, err := c.client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("erro ao enviar mensagem: %w", err)
	}

	fmt.Printf("[whatsapp] Mensagem enviada para %s\n", phone)
	return nil
}

// Desconectar disconnects the client.
func (c *Client) Desconectar() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		c.client.Disconnect()
		c.connected = false
	}
}

// IsConnected returns whether the client is connected.
func (c *Client) IsConnected() bool {
	return c.connected
}

// limparTelefone normalizes a phone number for WhatsApp.
func limparTelefone(tel string) string {
	// Remove everything that's not a digit
	var clean strings.Builder
	for _, r := range tel {
		if r >= '0' && r <= '9' {
			clean.WriteRune(r)
		}
	}
	phone := clean.String()

	if len(phone) < 10 {
		return ""
	}

	// If starts with 0, assume local (Brazil) and prepend 55
	if strings.HasPrefix(phone, "0") {
		phone = "55" + phone[1:]
	}

	// If doesn't start with country code, assume Brazil
	if !strings.HasPrefix(phone, "55") && len(phone) <= 11 {
		phone = "55" + phone
	}

	return phone
}

// printQR renders a basic QR representation for terminal.
func printQR(code string) {
	// Simple text-based QR display
	fmt.Println("┌─────────────────────────────────────┐")
	fmt.Println("│                                     │")
	fmt.Printf("│  QR Code: %-26s │\n", code[:20]+"...")
	fmt.Println("│                                     │")
	fmt.Println("│  Use seu celular para escanear:     │")
	fmt.Println("│  WhatsApp > Dispositivos Conectados │")
	fmt.Println("│  > Conectar Dispositivo             │")
	fmt.Println("│                                     │")
	fmt.Printf("│  Ou acesse: wa.me/qr/%-15s │\n", "")
	fmt.Println("│                                     │")
	fmt.Println("└─────────────────────────────────────┘")
	fmt.Println()
	// Also write the full code for QR generators
	fmt.Fprintf(os.Stderr, "[whatsapp] QR raw: %s\n", code)
}
