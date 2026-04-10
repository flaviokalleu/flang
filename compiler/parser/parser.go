package parser

import (
	"fmt"
	"strings"

	"github.com/flavio/flang/compiler/ast"
	"github.com/flavio/flang/compiler/lexer"
)

// Parser converts a stream of tokens into an AST.
type Parser struct {
	tokens  []lexer.Token
	pos     int
	program *ast.Program
}

// New creates a new Parser for the given tokens.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		pos:     0,
		program: &ast.Program{},
	}
}

// Parse processes all tokens and returns the AST.
func (p *Parser) Parse() (*ast.Program, error) {
	p.skipWhitespace()

	for !p.isAtEnd() {
		tok := p.current()
		switch tok.Type {
		case lexer.TokenImportar:
			if err := p.parseImportar(); err != nil {
				return nil, err
			}
		case lexer.TokenSistema:
			if err := p.parseSistema(); err != nil {
				return nil, err
			}
		case lexer.TokenDados:
			if err := p.parseDados(); err != nil {
				return nil, err
			}
		case lexer.TokenTelas:
			if err := p.parseTelas(); err != nil {
				return nil, err
			}
		case lexer.TokenEventos:
			if err := p.parseEventos(); err != nil {
				return nil, err
			}
		case lexer.TokenAcoes:
			if err := p.parseAcoes(); err != nil {
				return nil, err
			}
		case lexer.TokenTema:
			if err := p.parseTema(); err != nil {
				return nil, err
			}
		case lexer.TokenLogica:
			if err := p.parseLogica(); err != nil {
				return nil, err
			}
		case lexer.TokenBanco:
			if err := p.parseBanco(); err != nil {
				return nil, err
			}
		case lexer.TokenAutenticacao:
			if err := p.parseAuth(); err != nil {
				return nil, err
			}
		case lexer.TokenIntegracoes:
			if err := p.parseIntegracoes(); err != nil {
				return nil, err
			}
		default:
			p.advance()
		}
		p.skipWhitespace()
	}

	return p.program, nil
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) expect(tt lexer.TokenType) (lexer.Token, error) {
	tok := p.current()
	if tok.Type != tt {
		return tok, fmt.Errorf("line %d: expected token type %d, got %d (%q)", tok.Line, tt, tok.Type, tok.Value)
	}
	p.advance()
	return tok, nil
}

func (p *Parser) isAtEnd() bool {
	return p.pos >= len(p.tokens) || p.tokens[p.pos].Type == lexer.TokenEOF
}

func (p *Parser) skipNewlines() {
	for !p.isAtEnd() && p.current().Type == lexer.TokenNewline {
		p.advance()
	}
}

func (p *Parser) skipWhitespace() {
	for !p.isAtEnd() {
		tt := p.current().Type
		if tt == lexer.TokenNewline || tt == lexer.TokenIndent {
			p.advance()
		} else {
			break
		}
	}
}

func (p *Parser) skipToNextLine() {
	for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		p.advance()
	}
	p.skipNewlines()
}

func (p *Parser) isBlockKeyword() bool {
	return lexer.IsBlockKeyword(p.current().Type)
}

// parseSistema parses: sistema <name>
func (p *Parser) parseSistema() error {
	p.advance() // consume 'sistema'
	p.skipIndent()

	name := p.advance()
	if name.Type != lexer.TokenIdentifier && name.Type != lexer.TokenString {
		return fmt.Errorf("line %d: expected system name after 'sistema'", name.Line)
	}

	p.program.System = &ast.System{Name: name.Value}
	p.skipToNextLine()
	return nil
}

func (p *Parser) skipIndent() {
	for !p.isAtEnd() && p.current().Type == lexer.TokenIndent {
		p.advance()
	}
}

// isNameToken returns true if the token can be used as an identifier/name
// (includes actual identifiers and type keywords that may be used as field/model names).
func (p *Parser) isNameToken(tok lexer.Token) bool {
	// A name token is anything that's not whitespace, punctuation, or a block keyword
	switch tok.Type {
	case lexer.TokenEOF, lexer.TokenNewline, lexer.TokenIndent,
		lexer.TokenColon, lexer.TokenDot, lexer.TokenComma,
		lexer.TokenString, lexer.TokenNumber,
		lexer.TokenEquals, lexer.TokenPlus, lexer.TokenMinus,
		lexer.TokenStar, lexer.TokenSlash, lexer.TokenLParen, lexer.TokenRParen:
		return false
	}
	if lexer.IsBlockKeyword(tok.Type) {
		return false
	}
	return true
}

// parseDados parses the data models block.
func (p *Parser) parseDados() error {
	p.advance() // consume 'dados'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()
		// A model name is any name token where the next meaningful token is NOT a colon
		// (if it were a colon, it would be a field definition, not a model)
		if p.isNameToken(tok) && p.peekNextMeaningful().Type != lexer.TokenColon {
			model, err := p.parseModel()
			if err != nil {
				return err
			}
			p.program.Models = append(p.program.Models, model)
		} else {
			p.advance()
		}
	}
	return nil
}

// parseModel parses a single model definition.
func (p *Parser) parseModel() (*ast.Model, error) {
	nameTok := p.advance()
	model := &ast.Model{Name: nameTok.Value}

	// Check for modifiers after model name on the same line (e.g. soft_delete)
	for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		if p.current().Type == lexer.TokenSoftDelete {
			model.SoftDelete = true
			p.advance()
		} else if p.current().Type == lexer.TokenIndent {
			p.advance()
		} else {
			break
		}
	}

	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}
		if tok.Type == lexer.TokenNewline {
			p.advance()
			continue
		}

		// Any name token: check if it's a field (followed by colon) or a new model
		if p.isNameToken(tok) {
			if p.peekNextMeaningful().Type == lexer.TokenColon {
				// It's a field definition
				field, err := p.parseField()
				if err != nil {
					return nil, err
				}
				if field != nil {
					model.Fields = append(model.Fields, field)
				}
			} else {
				// It's a new model name — stop
				break
			}
		} else {
			p.advance()
		}
	}

	return model, nil
}

func (p *Parser) peekNextMeaningful() lexer.Token {
	saved := p.pos
	p.pos++
	for p.pos < len(p.tokens) {
		tok := p.tokens[p.pos]
		if tok.Type != lexer.TokenNewline && tok.Type != lexer.TokenIndent {
			p.pos = saved
			return tok
		}
		p.pos++
	}
	p.pos = saved
	return lexer.Token{Type: lexer.TokenEOF}
}

// parseField parses: <name>: <type>
func (p *Parser) parseField() (*ast.Field, error) {
	nameTok := p.advance() // consume field name (could be identifier or type keyword used as name)

	// Expect colon
	if p.current().Type != lexer.TokenColon {
		// Not a field definition
		p.pos--
		return nil, nil
	}
	p.advance() // consume ':'

	// Skip any indent tokens after colon
	p.skipIndent()

	typeTok := p.advance()
	fieldType, err := tokenToFieldType(typeTok)
	if err != nil {
		return nil, fmt.Errorf("line %d: %w", typeTok.Line, err)
	}

	field := &ast.Field{
		Name: nameTok.Value,
		Type: fieldType,
	}

	// Check for modifiers on the same line
	for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		switch p.current().Type {
		case lexer.TokenObrigatorio:
			field.Required = true
			p.advance()
		case lexer.TokenUnico:
			field.Unique = true
			p.advance()
		case lexer.TokenIndice:
			field.Index = true
			p.advance()
		case lexer.TokenPertenceA:
			p.advance()
			p.skipIndent()
			if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				field.Reference = p.advance().Value
			}
		default:
			p.advance()
		}
	}

	return field, nil
}

func tokenToFieldType(tok lexer.Token) (ast.FieldType, error) {
	switch tok.Type {
	case lexer.TokenTexto:
		return ast.FieldTexto, nil
	case lexer.TokenNumero:
		return ast.FieldNumero, nil
	case lexer.TokenData:
		return ast.FieldData, nil
	case lexer.TokenBooleano:
		return ast.FieldBooleano, nil
	case lexer.TokenEmail:
		return ast.FieldEmail, nil
	case lexer.TokenTelefone:
		return ast.FieldTelefone, nil
	case lexer.TokenImagem:
		return ast.FieldImagem, nil
	case lexer.TokenArquivo:
		return ast.FieldArquivo, nil
	case lexer.TokenUpload:
		return ast.FieldUpload, nil
	case lexer.TokenLink:
		return ast.FieldLink, nil
	case lexer.TokenStatus:
		return ast.FieldStatus, nil
	case lexer.TokenDinheiro:
		return ast.FieldDinheiro, nil
	case lexer.TokenSenha:
		return ast.FieldSenha, nil
	case lexer.TokenTextoLongo:
		return ast.FieldTextoLongo, nil
	case lexer.TokenEnum:
		return ast.FieldEnum, nil
	}
	typeMap := map[string]ast.FieldType{
		"texto": ast.FieldTexto, "numero": ast.FieldNumero, "data": ast.FieldData,
		"booleano": ast.FieldBooleano, "email": ast.FieldEmail, "telefone": ast.FieldTelefone,
		"imagem": ast.FieldImagem, "arquivo": ast.FieldArquivo, "upload": ast.FieldUpload,
		"link": ast.FieldLink, "status": ast.FieldStatus, "dinheiro": ast.FieldDinheiro,
		"senha": ast.FieldSenha, "texto_longo": ast.FieldTextoLongo, "enum": ast.FieldEnum,
		// EN
		"text": ast.FieldTexto, "number": ast.FieldNumero, "date": ast.FieldData,
		"boolean": ast.FieldBooleano, "phone": ast.FieldTelefone, "image": ast.FieldImagem,
		"file": ast.FieldArquivo, "money": ast.FieldDinheiro, "password": ast.FieldSenha,
		"long_text": ast.FieldTextoLongo,
	}
	if ft, ok := typeMap[tok.Value]; ok {
		return ft, nil
	}
	return "", fmt.Errorf("tipo desconhecido: %q", tok.Value)
}

// parseTelas parses the screens block.
func (p *Parser) parseTelas() error {
	p.advance() // consume 'telas'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()
		if tok.Type == lexer.TokenTela {
			screen, err := p.parseScreen()
			if err != nil {
				return err
			}
			p.program.Screens = append(p.program.Screens, screen)
		} else {
			p.advance()
		}
	}
	return nil
}

// parseScreen parses a single screen definition.
func (p *Parser) parseScreen() (*ast.Screen, error) {
	p.advance() // consume 'tela'
	p.skipIndent()

	nameTok := p.advance()
	screen := &ast.Screen{Name: nameTok.Value}
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		// Check if we hit a new 'tela' keyword
		if tok.Type == lexer.TokenTela {
			break
		}

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		switch tok.Type {
		case lexer.TokenTitulo:
			p.advance()
			p.skipIndent()
			if p.current().Type == lexer.TokenString {
				screen.Title = p.advance().Value
			}
		case lexer.TokenLista:
			comp, err := p.parseListComponent()
			if err != nil {
				return nil, err
			}
			screen.Components = append(screen.Components, comp)
		case lexer.TokenBotao:
			comp, err := p.parseButtonComponent()
			if err != nil {
				return nil, err
			}
			screen.Components = append(screen.Components, comp)
		case lexer.TokenFormulario:
			comp, err := p.parseFormComponent()
			if err != nil {
				return nil, err
			}
			screen.Components = append(screen.Components, comp)
		case lexer.TokenMostrar:
			comp := p.parseShowComponent()
			screen.Components = append(screen.Components, comp)
		default:
			p.advance()
		}
	}

	return screen, nil
}

func (p *Parser) parseListComponent() (*ast.Component, error) {
	p.advance() // consume 'lista'
	p.skipIndent()

	comp := &ast.Component{
		Type:       ast.CompList,
		Properties: make(map[string]string),
	}

	if !p.isAtEnd() && (p.current().Type == lexer.TokenIdentifier || lexer.IsTypeKeyword(p.current().Type)) {
		comp.Target = p.advance().Value
	}

	p.skipWhitespace()

	// Parse children (mostrar fields)
	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()
		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}
		if tok.Type == lexer.TokenMostrar {
			child := p.parseShowComponent()
			comp.Children = append(comp.Children, child)
			continue
		}
		// Any other token means we're done with list children
		break
	}

	return comp, nil
}

func (p *Parser) parseShowComponent() *ast.Component {
	p.advance() // consume 'mostrar'
	p.skipIndent()

	comp := &ast.Component{
		Type:       ast.CompShow,
		Properties: make(map[string]string),
	}

	if !p.isAtEnd() {
		tok := p.current()
		if tok.Type == lexer.TokenIdentifier || lexer.IsTypeKeyword(tok.Type) {
			comp.Target = p.advance().Value
		}
	}

	return comp
}

func (p *Parser) parseButtonComponent() (*ast.Component, error) {
	p.advance() // consume 'botao'
	p.skipIndent()

	comp := &ast.Component{
		Type:       ast.CompButton,
		Properties: make(map[string]string),
	}

	// Optional color
	if p.current().Type == lexer.TokenIdentifier {
		comp.Properties["color"] = p.advance().Value
	}

	p.skipWhitespace()

	// Parse button properties
	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()
		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}
		if tok.Type == lexer.TokenTexto || tok.Value == "texto" {
			p.advance()
			p.skipIndent()
			if p.current().Type == lexer.TokenString {
				comp.Properties["text"] = p.advance().Value
			}
			continue
		}
		break
	}

	return comp, nil
}

func (p *Parser) parseFormComponent() (*ast.Component, error) {
	p.advance() // consume 'formulario'
	p.skipIndent()

	comp := &ast.Component{
		Type:       ast.CompForm,
		Properties: make(map[string]string),
	}

	if !p.isAtEnd() && p.current().Type == lexer.TokenIdentifier {
		comp.Target = p.advance().Value
	}

	return comp, nil
}

// parseEventos parses the events block.
func (p *Parser) parseEventos() error {
	p.advance() // consume 'eventos'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()
		if tok.Type == lexer.TokenQuando {
			event, err := p.parseEvent()
			if err != nil {
				return err
			}
			p.program.Events = append(p.program.Events, event)
		} else {
			p.advance()
		}
	}
	return nil
}

func (p *Parser) parseEvent() (*ast.Event, error) {
	p.advance() // consume 'quando'
	p.skipIndent()

	event := &ast.Event{}

	// Parse trigger: clicar, enviar, etc.
	if !p.isAtEnd() {
		event.Trigger = p.advance().Value
	}

	// Parse target (quoted string)
	p.skipIndent()
	if p.current().Type == lexer.TokenString {
		event.Target = p.advance().Value
	}

	p.skipWhitespace()

	// Parse action reference (next line typically)
	var actionParts []string
	for !p.isAtEnd() && p.current().Type != lexer.TokenNewline && !p.isBlockKeyword() && p.current().Type != lexer.TokenQuando {
		if p.current().Type == lexer.TokenIndent {
			p.advance()
			continue
		}
		actionParts = append(actionParts, p.advance().Value)
	}
	event.ActionRef = strings.Join(actionParts, " ")

	return event, nil
}

// parseAcoes parses the actions block.
func (p *Parser) parseAcoes() error {
	p.advance() // consume 'acoes'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		p.advance()
	}
	return nil
}

// parseTema parses the theme block.
func (p *Parser) parseTema() error {
	p.advance() // consume 'tema'
	p.skipWhitespace()

	theme := ast.DefaultTheme()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// Parse key: value pairs or keywords
		switch {
		case tok.Value == "cor" || tok.Type == lexer.TokenCor:
			p.advance()
			p.skipIndent()
			// cor primaria "#hex" / cor secundaria "#hex" etc
			if p.current().Type == lexer.TokenIdentifier || p.current().Type == lexer.TokenCor {
				which := p.advance().Value
				p.skipIndent()
				if p.current().Type == lexer.TokenString {
					val := p.advance().Value
					switch which {
					case "primaria", "primary":
						theme.Primary = val
					case "secundaria", "secondary":
						theme.Secondary = val
					case "destaque", "accent":
						theme.Accent = val
					case "sidebar":
						theme.Sidebar = val
					}
				}
			} else if p.current().Type == lexer.TokenString {
				theme.Primary = p.advance().Value
			}
		case tok.Value == "escuro" || tok.Type == lexer.TokenEscuro:
			p.advance()
			theme.Dark = true
		case tok.Value == "icone" || tok.Type == lexer.TokenIcone:
			p.advance()
			p.skipIndent()
			if p.current().Type == lexer.TokenString {
				theme.Icon = p.advance().Value
			}
		default:
			p.advance()
		}
	}

	p.program.Theme = theme
	return nil
}

// parseImportar parses: importar "arquivo.fg"
//                       importar dados de "arquivo.fg"
//                       importar tela de "arquivo.fg"
func (p *Parser) parseImportar() error {
	p.advance() // consume 'importar'
	p.skipIndent()

	imp := &ast.Import{}

	tok := p.current()

	// importar "file.fg" (import everything)
	if tok.Type == lexer.TokenString {
		imp.What = "tudo"
		imp.Path = p.advance().Value
		p.program.Imports = append(p.program.Imports, imp)
		return nil
	}

	// importar <what> de "file.fg"
	imp.What = p.advance().Value
	p.skipIndent()

	// expect 'de'
	if p.current().Type == lexer.TokenDe {
		p.advance()
		p.skipIndent()
	}

	if p.current().Type == lexer.TokenString {
		imp.Path = p.advance().Value
	} else {
		return fmt.Errorf("line %d: esperado caminho do arquivo após 'importar %s de'", p.current().Line, imp.What)
	}

	p.program.Imports = append(p.program.Imports, imp)
	return nil
}

// parseLogica parses the logic block with rules.
func (p *Parser) parseLogica() error {
	p.advance() // consume 'logica'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// se <field> <operator> <value>
		//   <action> <arg>
		if tok.Type == lexer.TokenSe || tok.Type == lexer.TokenQuando {
			if err := p.parseRule(); err != nil {
				return err
			}
			continue
		}

		// validar <model> <field> <condition>
		if tok.Type == lexer.TokenValidar {
			if err := p.parseValidacao(); err != nil {
				return err
			}
			continue
		}

		p.advance()
	}
	return nil
}

// parseRule parses: se <field> igual/maior/menor <value> \n <action> <arg>
func (p *Parser) parseRule() error {
	p.advance() // consume 'se' or 'quando'
	p.skipIndent()

	rule := &ast.Rule{}

	// field name
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		rule.Field = p.advance().Value
	}
	p.skipIndent()

	// operator
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		rule.Operator = p.advance().Value
	}
	p.skipIndent()

	// value
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		if p.current().Type == lexer.TokenString {
			rule.Value = p.advance().Value
		} else {
			rule.Value = p.advance().Value
		}
	}

	p.skipWhitespace()

	// action line: mudar/validar/calcular/definir <arg>
	if !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()
		if tok.Type != lexer.TokenSe && tok.Type != lexer.TokenQuando && tok.Type != lexer.TokenValidar {
			rule.Action = p.advance().Value
			p.skipIndent()

			// Collect rest of line as arg
			var parts []string
			for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				if p.current().Type == lexer.TokenIndent {
					p.advance()
					continue
				}
				if p.current().Type == lexer.TokenString {
					parts = append(parts, p.advance().Value)
				} else {
					parts = append(parts, p.advance().Value)
				}
			}
			rule.ActionArg = strings.Join(parts, " ")
		}
	}

	p.program.Rules = append(p.program.Rules, rule)
	return nil
}

// parseValidacao parses: validar <field> <condition>
func (p *Parser) parseValidacao() error {
	p.advance() // consume 'validar'
	p.skipIndent()

	rule := &ast.Rule{Action: "validar"}

	// field
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		rule.Field = p.advance().Value
	}
	p.skipIndent()

	// operator
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		rule.Operator = p.advance().Value
	}
	p.skipIndent()

	// value
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		if p.current().Type == lexer.TokenString {
			rule.Value = p.advance().Value
		} else {
			rule.Value = p.advance().Value
		}
	}

	p.program.Rules = append(p.program.Rules, rule)
	return nil
}

// parseBanco parses database configuration block.
func (p *Parser) parseBanco() error {
	p.advance() // consume 'banco'
	p.skipWhitespace()

	db := ast.DefaultDatabase()

	// Check if driver name is on same line: banco postgres
	if !p.isAtEnd() && p.current().Type == lexer.TokenIdentifier {
		db.Driver = p.advance().Value
	}
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		if tok.Type == lexer.TokenIdentifier || lexer.IsTypeKeyword(tok.Type) {
			key := p.advance().Value
			p.skipIndent()

			if p.current().Type == lexer.TokenColon {
				p.advance()
				p.skipIndent()
			}

			val := ""
			if p.current().Type == lexer.TokenString {
				val = p.advance().Value
			} else if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				val = p.advance().Value
			}

			switch key {
			case "driver", "tipo", "type":
				db.Driver = val
			case "host", "servidor", "server":
				db.Host = val
			case "port", "porta":
				db.Port = val
			case "nome", "name", "database", "db":
				db.Name = val
			case "usuario", "user", "username":
				db.User = val
			case "senha", "password", "pass":
				db.Password = val
			}
			continue
		}

		p.advance()
	}

	p.program.Database = db
	return nil
}

// parseIntegracoes parses the integrations block (whatsapp, etc).
func (p *Parser) parseIntegracoes() error {
	p.advance() // consume 'integracoes'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		if tok.Type == lexer.TokenWhatsapp {
			if err := p.parseWhatsApp(); err != nil {
				return err
			}
			continue
		}

		if tok.Type == lexer.TokenEmailInteg || tok.Value == "email" {
			if err := p.parseEmailInteg(); err != nil {
				return err
			}
			continue
		}

		if tok.Type == lexer.TokenCron {
			if err := p.parseCron(); err != nil {
				return err
			}
			continue
		}

		p.advance()
	}
	return nil
}

// parseWhatsApp parses the whatsapp sub-block inside integracoes.
// whatsapp
//   quando criar pedido
//     enviar mensagem para cliente.telefone
//       texto "Seu pedido foi recebido!"
func (p *Parser) parseWhatsApp() error {
	p.advance() // consume 'whatsapp'
	p.skipWhitespace()

	// Enable WhatsApp
	if p.program.WhatsApp == nil {
		p.program.WhatsApp = &ast.WhatsAppConfig{Enabled: true, DBPath: "whatsapp.db"}
	}

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// quando <trigger> <model>
		if tok.Type == lexer.TokenQuando {
			notif, err := p.parseNotifier("whatsapp")
			if err != nil {
				return err
			}
			if notif != nil {
				p.program.Notifiers = append(p.program.Notifiers, notif)
			}
			continue
		}

		// Exit if we hit another integration section
		if tok.Type == lexer.TokenWhatsapp {
			break
		}

		p.advance()
	}

	return nil
}

// parseNotifier parses a notification trigger.
// quando criar pedido
//   enviar mensagem para cliente.telefone
//     texto "Mensagem aqui"
func (p *Parser) parseNotifier(channel string) (*ast.Notifier, error) {
	p.advance() // consume 'quando'
	p.skipIndent()

	notif := &ast.Notifier{Channel: channel}

	// trigger: criar, atualizar, deletar, or field condition
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		triggerTok := p.advance()
		notif.Trigger = triggerTok.Value
	}
	p.skipIndent()

	// model name or field condition
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		notif.Model = p.advance().Value
	}

	// Optional condition: igual "value"
	p.skipIndent()
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		if p.current().Type == lexer.TokenIgual || p.current().Value == "igual" || p.current().Value == "equals" {
			p.advance() // skip 'igual'
			p.skipIndent()
			// The model was actually the field, trigger was condition context
			notif.Field = notif.Model
			notif.Model = ""
			if p.current().Type == lexer.TokenString {
				notif.Value = p.advance().Value
			} else if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				notif.Value = p.advance().Value
			}
		}
	}

	p.skipWhitespace()

	// Parse body lines: enviar mensagem para <dest> / texto "..."
	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// Stop at next 'quando'
		if tok.Type == lexer.TokenQuando {
			break
		}

		switch tok.Value {
		case "enviar", "send":
			p.advance()
			p.skipIndent()
			// skip 'mensagem' / 'message'
			if p.current().Type == lexer.TokenMensagem {
				p.advance()
			}
			p.skipIndent()
			// skip 'para' / 'to'
			if p.current().Value == "para" || p.current().Type == lexer.TokenPara {
				p.advance()
			}
			p.skipIndent()
			// destination
			if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				notif.SendTo = p.advance().Value
				// Check for dotted access like cliente.telefone
				if p.current().Type == lexer.TokenDot {
					p.advance()
					if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
						notif.SendTo += "." + p.advance().Value
					}
				}
			}

		case "texto", "text":
			p.advance()
			p.skipIndent()
			if p.current().Type == lexer.TokenString {
				notif.Message = p.advance().Value
			}

		default:
			p.advance()
		}
	}

	return notif, nil
}

// parseAuth parses the authentication block.
// autenticacao
//   usuario: email
//   senha: senha
//   roles: admin, usuario
func (p *Parser) parseAuth() error {
	p.advance() // consume 'autenticacao'
	p.skipWhitespace()

	auth := &ast.AuthConfig{
		Enabled:    true,
		UserModel:  "usuario",
		LoginField: "email",
		PassField:  "senha",
		JWTSecret:  "flang-secret-change-me",
	}

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// Accept any token as key (keywords can be used as config keys)
		if tok.Type != lexer.TokenNewline && tok.Type != lexer.TokenIndent && tok.Type != lexer.TokenEOF {
			key := p.advance().Value
			p.skipIndent()

			if p.current().Type == lexer.TokenColon {
				p.advance()
				p.skipIndent()
			}

			val := ""
			if p.current().Type == lexer.TokenString {
				val = p.advance().Value
			} else if !p.isAtEnd() && p.current().Type != lexer.TokenNewline && p.current().Type != lexer.TokenComma {
				val = p.advance().Value
			}

			switch key {
			case "usuario", "user", "modelo", "model":
				auth.UserModel = val
			case "login", "campo_login", "login_field":
				auth.LoginField = val
			case "senha", "password", "campo_senha", "password_field":
				auth.PassField = val
			case "secret", "segredo", "jwt_secret":
				auth.JWTSecret = val
			case "roles", "permissoes":
				auth.Roles = append(auth.Roles, val)
				for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
					if p.current().Type == lexer.TokenComma || p.current().Type == lexer.TokenIndent {
						p.advance()
						continue
					}
					auth.Roles = append(auth.Roles, p.advance().Value)
				}
			}
			continue
		}

		p.advance()
	}

	p.program.Auth = auth
	return nil
}

// parseEmailInteg parses the email sub-block inside integracoes.
// email
//
//	servidor: "smtp.gmail.com"
//	porta: "587"
//	usuario: "me@gmail.com"
//	senha: "app-password"
//	quando criar pedido
//	  enviar email para cliente.email
//	    assunto "Pedido recebido"
//	    texto "Olá {cliente}, seu pedido..."
func (p *Parser) parseEmailInteg() error {
	p.advance() // consume 'email'
	p.skipWhitespace()

	// Initialize email config
	if p.program.Email == nil {
		p.program.Email = &ast.EmailConfig{}
	}

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// quando <trigger> <model> — parse as email notifier
		if tok.Type == lexer.TokenQuando {
			notif, err := p.parseEmailNotifier()
			if err != nil {
				return err
			}
			if notif != nil {
				p.program.Notifiers = append(p.program.Notifiers, notif)
			}
			continue
		}

		// Exit if we hit another integration keyword
		if tok.Type == lexer.TokenWhatsapp || tok.Type == lexer.TokenCron || tok.Type == lexer.TokenEmailInteg {
			break
		}

		// Config key: value pairs (servidor, porta, usuario, senha)
		if tok.Type == lexer.TokenIdentifier || lexer.IsTypeKeyword(tok.Type) ||
			tok.Type == lexer.TokenSenha || tok.Type == lexer.TokenUsuario {
			key := p.advance().Value
			p.skipIndent()

			if p.current().Type == lexer.TokenColon {
				p.advance()
				p.skipIndent()
			}

			val := ""
			if p.current().Type == lexer.TokenString {
				val = p.advance().Value
			} else if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				val = p.advance().Value
			}

			switch key {
			case "servidor", "server", "host":
				p.program.Email.Host = val
			case "porta", "port":
				p.program.Email.Port = val
			case "usuario", "user", "username":
				p.program.Email.User = val
			case "senha", "password", "pass":
				p.program.Email.Password = val
			case "de", "from", "remetente":
				p.program.Email.From = val
			}
			continue
		}

		p.advance()
	}

	return nil
}

// parseEmailNotifier parses an email notification trigger.
// quando criar pedido
//
//	enviar email para cliente.email
//	  assunto "Pedido recebido"
//	  texto "Mensagem aqui"
func (p *Parser) parseEmailNotifier() (*ast.Notifier, error) {
	p.advance() // consume 'quando'
	p.skipIndent()

	notif := &ast.Notifier{Channel: "email"}

	// trigger: criar, atualizar, deletar
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		notif.Trigger = p.advance().Value
	}
	p.skipIndent()

	// model name
	if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		notif.Model = p.advance().Value
	}

	p.skipWhitespace()

	// Parse body: enviar email para <dest> / assunto "..." / texto "..."
	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// Stop at next 'quando' or integration keyword
		if tok.Type == lexer.TokenQuando || tok.Type == lexer.TokenWhatsapp ||
			tok.Type == lexer.TokenCron || tok.Type == lexer.TokenEmailInteg {
			break
		}

		switch tok.Value {
		case "enviar", "send":
			p.advance()
			p.skipIndent()
			// skip 'email' / 'mensagem' / 'message'
			if p.current().Value == "email" || p.current().Type == lexer.TokenMensagem {
				p.advance()
			}
			p.skipIndent()
			// skip 'para' / 'to'
			if p.current().Value == "para" || p.current().Type == lexer.TokenPara {
				p.advance()
			}
			p.skipIndent()
			// destination
			if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
				notif.SendTo = p.advance().Value
				if p.current().Type == lexer.TokenDot {
					p.advance()
					if !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
						notif.SendTo += "." + p.advance().Value
					}
				}
			}

		case "assunto", "subject":
			p.advance()
			p.skipIndent()
			if p.current().Type == lexer.TokenString {
				notif.Subject = p.advance().Value
			}

		case "texto", "text":
			p.advance()
			p.skipIndent()
			if p.current().Type == lexer.TokenString {
				notif.Message = p.advance().Value
			}

		default:
			p.advance()
		}
	}

	return notif, nil
}

// parseCron parses the cron sub-block inside integracoes.
// cron
//
//	cada 5 minutos
//	  chamar api "https://example.com/webhook"
//	cada 1 hora
//	  limpar sessoes
func (p *Parser) parseCron() error {
	p.advance() // consume 'cron'
	p.skipWhitespace()

	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// cada/every <N> <unit>
		if tok.Type == lexer.TokenCada {
			job, err := p.parseCronJob()
			if err != nil {
				return err
			}
			if job != nil {
				p.program.Crons = append(p.program.Crons, job)
			}
			continue
		}

		// Exit if we hit another integration keyword
		if tok.Type == lexer.TokenWhatsapp || tok.Type == lexer.TokenEmailInteg || tok.Type == lexer.TokenCron {
			break
		}

		p.advance()
	}

	return nil
}

// parseCronJob parses a single cron job definition.
// cada 5 minutos
//
//	chamar api "https://example.com/webhook"
func (p *Parser) parseCronJob() (*ast.CronJob, error) {
	p.advance() // consume 'cada' / 'every'
	p.skipIndent()

	job := &ast.CronJob{}

	// Parse interval: <number> <unit>
	var intervalParts []string
	for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
		if p.current().Type == lexer.TokenIndent {
			p.advance()
			continue
		}
		intervalParts = append(intervalParts, p.advance().Value)
	}
	job.Every = strings.Join(intervalParts, " ")

	p.skipWhitespace()

	// Parse action line(s)
	for !p.isAtEnd() && !p.isBlockKeyword() {
		tok := p.current()

		if tok.Type == lexer.TokenNewline || tok.Type == lexer.TokenIndent {
			p.advance()
			continue
		}

		// Stop at next 'cada' or integration keyword
		if tok.Type == lexer.TokenCada || tok.Type == lexer.TokenWhatsapp ||
			tok.Type == lexer.TokenEmailInteg || tok.Type == lexer.TokenCron {
			break
		}

		// chamar api "url" / chamar "url"
		if tok.Type == lexer.TokenChamar || tok.Value == "chamar" || tok.Value == "call" {
			p.advance()
			p.skipIndent()
			job.Action = "chamar"

			// skip optional 'api'
			if p.current().Type == lexer.TokenApi || p.current().Value == "api" {
				p.advance()
				p.skipIndent()
			}

			// URL
			if p.current().Type == lexer.TokenString {
				job.Target = p.advance().Value
			}
			break
		}

		// Generic action: collect remaining tokens on the line as action + target
		var actionParts []string
		for !p.isAtEnd() && p.current().Type != lexer.TokenNewline {
			if p.current().Type == lexer.TokenIndent {
				p.advance()
				continue
			}
			actionParts = append(actionParts, p.advance().Value)
		}
		if len(actionParts) > 0 {
			job.Action = actionParts[0]
			if len(actionParts) > 1 {
				job.Target = strings.Join(actionParts[1:], " ")
			}
		}
		break
	}

	return job, nil
}
