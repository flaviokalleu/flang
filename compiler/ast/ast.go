package ast

// Node is the base interface for all AST nodes.
type Node interface {
	NodeType() string
}

// Program is the root AST node representing a complete .fg file.
type Program struct {
	System    *System
	Theme     *Theme
	Database  *DatabaseConfig
	WhatsApp  *WhatsAppConfig
	Imports   []*Import
	Models    []*Model
	Screens   []*Screen
	Events    []*Event
	Actions   []*Action
	Rules     []*Rule
	Notifiers []*Notifier
}

// WhatsAppConfig holds WhatsApp connection settings.
type WhatsAppConfig struct {
	Enabled bool
	DBPath  string // path to whatsmeow session store
}

func (w *WhatsAppConfig) NodeType() string { return "WhatsAppConfig" }

// Notifier represents a notification trigger.
// quando criar pedido → enviar mensagem para cliente.telefone texto "..."
type Notifier struct {
	Trigger  string // "criar", "atualizar", "deletar"
	Model    string // model name
	Field    string // condition field (e.g. "status")
	Value    string // condition value (e.g. "pronto")
	SendTo   string // destination field (e.g. "cliente.telefone" or phone number)
	Message  string // message template
	Channel  string // "whatsapp", "email", etc
}

func (n *Notifier) NodeType() string { return "Notifier" }

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Driver   string // sqlite, mysql, postgres
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (d *DatabaseConfig) NodeType() string { return "DatabaseConfig" }

// DefaultDatabase returns SQLite config.
func DefaultDatabase() *DatabaseConfig {
	return &DatabaseConfig{Driver: "sqlite"}
}

func (p *Program) NodeType() string { return "Program" }

// Merge combines another program's definitions into this one (for imports).
func (p *Program) Merge(other *Program) {
	if other.Theme != nil && p.Theme == nil {
		p.Theme = other.Theme
	}
	p.Models = append(p.Models, other.Models...)
	p.Screens = append(p.Screens, other.Screens...)
	p.Events = append(p.Events, other.Events...)
	p.Actions = append(p.Actions, other.Actions...)
	p.Rules = append(p.Rules, other.Rules...)
}

// Import represents an import statement.
type Import struct {
	What string // what to import: "dados", "tela", "tudo", or specific name
	Path string // file path
}

func (i *Import) NodeType() string { return "Import" }

// System defines the application name and metadata.
type System struct {
	Name string
}

func (s *System) NodeType() string { return "System" }

// Theme holds visual customization.
type Theme struct {
	Primary   string
	Secondary string
	Accent    string
	Dark      bool
	Sidebar   string
	Icon      string
}

func (t *Theme) NodeType() string { return "Theme" }

// DefaultTheme returns the default theme.
func DefaultTheme() *Theme {
	return &Theme{
		Primary:   "#6366f1",
		Secondary: "#8b5cf6",
		Accent:    "#f59e0b",
		Dark:      false,
		Sidebar:   "#1e1b4b",
	}
}

// Model represents a data model.
type Model struct {
	Name   string
	Icon   string
	Fields []*Field
}

func (m *Model) NodeType() string { return "Model" }

// FieldType represents the Flang data types.
type FieldType string

const (
	FieldTexto    FieldType = "texto"
	FieldNumero   FieldType = "numero"
	FieldData     FieldType = "data"
	FieldBooleano FieldType = "booleano"
	FieldEmail    FieldType = "email"
	FieldTelefone FieldType = "telefone"
	FieldImagem   FieldType = "imagem"
	FieldArquivo  FieldType = "arquivo"
	FieldUpload   FieldType = "upload"
	FieldLink     FieldType = "link"
	FieldStatus   FieldType = "status"
	FieldDinheiro FieldType = "dinheiro"
	FieldSenha    FieldType = "senha"
)

// Field represents a field inside a model.
type Field struct {
	Name      string
	Type      FieldType
	Required  bool
	Unique    bool
	Default   string
	Reference string // pertence_a <model>
}

func (f *Field) NodeType() string { return "Field" }

// Screen represents a UI screen definition.
type Screen struct {
	Name       string
	Title      string
	Components []*Component
}

func (s *Screen) NodeType() string { return "Screen" }

// ComponentType identifies the kind of UI component.
type ComponentType string

const (
	CompList   ComponentType = "lista"
	CompShow   ComponentType = "mostrar"
	CompButton ComponentType = "botao"
	CompForm   ComponentType = "formulario"
	CompInput  ComponentType = "entrada"
	CompImage  ComponentType = "imagem"
	CompText   ComponentType = "texto"
	CompSearch ComponentType = "busca"
)

// Component represents a UI element inside a screen.
type Component struct {
	Type       ComponentType
	Target     string
	Properties map[string]string
	Children   []*Component
}

func (c *Component) NodeType() string { return "Component" }

// Event represents an event handler.
type Event struct {
	Trigger   string
	Target    string
	ActionRef string
}

func (e *Event) NodeType() string { return "Event" }

// Action represents a custom action block.
type Action struct {
	Name  string
	Steps []*ActionStep
}

func (a *Action) NodeType() string { return "Action" }

// ActionStep is a single step within an action.
type ActionStep struct {
	Command string
	Args    []string
}

func (s *ActionStep) NodeType() string { return "ActionStep" }

// Rule represents a logic rule (se/quando condition).
type Rule struct {
	Field     string // campo
	Operator  string // igual, maior, menor
	Value     string // valor para comparar
	Action    string // ação: mudar, validar, calcular, definir
	ActionArg string // argumento da ação
}

func (r *Rule) NodeType() string { return "Rule" }

// SQLType returns the SQLite type for a Flang field type.
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
