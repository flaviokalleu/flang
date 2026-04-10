package ast

// Node is the base interface for all AST nodes.
type Node interface {
	NodeType() string
}

// Program is the root AST node.
type Program struct {
	System    *System
	Theme     *Theme
	Database  *DatabaseConfig
	Auth      *AuthConfig
	WhatsApp  *WhatsAppConfig
	Email     *EmailConfig
	Imports   []*Import
	Models    []*Model
	Screens   []*Screen
	Events    []*Event
	Actions   []*Action
	Rules     []*Rule
	Notifiers []*Notifier
	Crons     []*CronJob
	Env       map[string]string
}

func (p *Program) NodeType() string { return "Program" }

// Merge combines another program into this one (for imports).
func (p *Program) Merge(other *Program) {
	if other.Theme != nil && p.Theme == nil {
		p.Theme = other.Theme
	}
	if other.Auth != nil && p.Auth == nil {
		p.Auth = other.Auth
	}
	if other.Email != nil && p.Email == nil {
		p.Email = other.Email
	}
	p.Models = append(p.Models, other.Models...)
	p.Screens = append(p.Screens, other.Screens...)
	p.Events = append(p.Events, other.Events...)
	p.Actions = append(p.Actions, other.Actions...)
	p.Rules = append(p.Rules, other.Rules...)
	p.Notifiers = append(p.Notifiers, other.Notifiers...)
	p.Crons = append(p.Crons, other.Crons...)
}

// ==================== System ====================

type System struct {
	Name string
}

func (s *System) NodeType() string { return "System" }

// ==================== Import ====================

type Import struct {
	What string
	Path string
}

func (i *Import) NodeType() string { return "Import" }

// ==================== Theme ====================

type Theme struct {
	Primary   string
	Secondary string
	Accent    string
	Dark      bool
	Sidebar   string
	Icon      string
}

func (t *Theme) NodeType() string { return "Theme" }

func DefaultTheme() *Theme {
	return &Theme{
		Primary: "#6366f1", Secondary: "#8b5cf6",
		Accent: "#f59e0b", Sidebar: "#1e1b4b",
	}
}

// ==================== Database ====================

type DatabaseConfig struct {
	Driver   string // sqlite, mysql, postgres
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (d *DatabaseConfig) NodeType() string { return "DatabaseConfig" }

func DefaultDatabase() *DatabaseConfig {
	return &DatabaseConfig{Driver: "sqlite"}
}

// ==================== Auth ====================

type AuthConfig struct {
	Enabled    bool
	UserModel  string   // model name for users (default: "usuario")
	LoginField string   // field used for login (default: "email")
	PassField  string   // password field (default: "senha")
	Roles      []string // available roles
	JWTSecret  string
}

func (a *AuthConfig) NodeType() string { return "AuthConfig" }

// ==================== WhatsApp ====================

type WhatsAppConfig struct {
	Enabled bool
	DBPath  string
}

func (w *WhatsAppConfig) NodeType() string { return "WhatsAppConfig" }

// ==================== Email ====================

type EmailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

func (e *EmailConfig) NodeType() string { return "EmailConfig" }

// ==================== Model ====================

type Model struct {
	Name       string
	Icon       string
	Fields     []*Field
	SoftDelete bool
	IsAuth     bool // is this the auth user model?
}

func (m *Model) NodeType() string { return "Model" }

// ==================== Field ====================

type FieldType string

const (
	FieldTexto     FieldType = "texto"
	FieldNumero    FieldType = "numero"
	FieldData      FieldType = "data"
	FieldBooleano  FieldType = "booleano"
	FieldEmail     FieldType = "email"
	FieldTelefone  FieldType = "telefone"
	FieldImagem    FieldType = "imagem"
	FieldArquivo   FieldType = "arquivo"
	FieldUpload    FieldType = "upload"
	FieldLink      FieldType = "link"
	FieldStatus    FieldType = "status"
	FieldDinheiro  FieldType = "dinheiro"
	FieldSenha     FieldType = "senha"
	FieldTextoLongo FieldType = "texto_longo"
	FieldEnum      FieldType = "enum"
)

type Field struct {
	Name       string
	Type       FieldType
	Required   bool
	Unique     bool
	Default    string
	Reference  string   // pertence_a model
	EnumValues []string // for enum type
	Index      bool
}

func (f *Field) NodeType() string { return "Field" }

func (ft FieldType) SQLType() string {
	switch ft {
	case FieldNumero, FieldDinheiro:
		return "REAL"
	case FieldBooleano:
		return "INTEGER"
	case FieldData:
		return "DATETIME"
	default:
		return "TEXT"
	}
}

// ==================== Screen ====================

type Screen struct {
	Name       string
	Title      string
	Public     bool // accessible without login
	Requires   string // required role
	Components []*Component
}

func (s *Screen) NodeType() string { return "Screen" }

// ==================== Component ====================

type ComponentType string

const (
	CompList     ComponentType = "lista"
	CompShow     ComponentType = "mostrar"
	CompButton   ComponentType = "botao"
	CompForm     ComponentType = "formulario"
	CompInput    ComponentType = "entrada"
	CompImage    ComponentType = "imagem"
	CompText     ComponentType = "texto"
	CompSearch   ComponentType = "busca"
	CompChart    ComponentType = "grafico"
	CompSelect   ComponentType = "selecionar"
	CompTextarea ComponentType = "area_texto"
)

type Component struct {
	Type       ComponentType
	Target     string
	Properties map[string]string
	Children   []*Component
}

func (c *Component) NodeType() string { return "Component" }

// ==================== Event ====================

type Event struct {
	Trigger   string
	Target    string
	ActionRef string
}

func (e *Event) NodeType() string { return "Event" }

// ==================== Action ====================

type Action struct {
	Name  string
	Steps []*ActionStep
}

func (a *Action) NodeType() string { return "Action" }

type ActionStep struct {
	Command string
	Args    []string
}

func (s *ActionStep) NodeType() string { return "ActionStep" }

// ==================== Rule ====================

type Rule struct {
	Field     string
	Operator  string
	Value     string
	Action    string
	ActionArg string
}

func (r *Rule) NodeType() string { return "Rule" }

// ==================== Notifier ====================

type Notifier struct {
	Trigger  string
	Model    string
	Field    string
	Value    string
	SendTo   string
	Message  string
	Subject  string // email subject
	Channel  string // whatsapp, email, webhook
}

func (n *Notifier) NodeType() string { return "Notifier" }

// ==================== CronJob ====================

type CronJob struct {
	Every    string // "1 hora", "30 minutos", etc
	Action   string // what to do
	Target   string // model or URL
}

func (c *CronJob) NodeType() string { return "CronJob" }
