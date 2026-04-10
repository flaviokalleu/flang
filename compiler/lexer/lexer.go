package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType identifies the kind of token.
type TokenType int

const (
	// Special
	TokenEOF TokenType = iota
	TokenNewline
	TokenIndent

	// Literals
	TokenIdentifier
	TokenString
	TokenNumber

	// Punctuation
	TokenColon
	TokenDot

	// Block keywords
	TokenSistema
	TokenDados
	TokenTelas
	TokenAcoes
	TokenEventos
	TokenIntegracoes
	TokenTema
	TokenLogica
	TokenBanco

	// Import
	TokenImportar
	TokenDe

	// Screen keywords
	TokenTela
	TokenTitulo
	TokenLista
	TokenMostrar
	TokenBotao
	TokenFormulario
	TokenEntrada
	TokenTextoKW
	TokenBusca
	TokenDashboard

	// Type keywords
	TokenTexto
	TokenNumero
	TokenData
	TokenBooleano
	TokenEmail
	TokenTelefone
	TokenImagem
	TokenArquivo
	TokenUpload
	TokenLink
	TokenStatus
	TokenDinheiro
	TokenSenha

	// Event keywords
	TokenQuando
	TokenClicar
	TokenCriar
	TokenAtualizar
	TokenDeletar
	TokenEnviar

	// Logic keywords
	TokenSe
	TokenSenao
	TokenIgual
	TokenMaior
	TokenMenor
	TokenE
	TokenOu
	TokenEntao
	TokenValidar
	TokenCalcular
	TokenDefinir
	TokenRetornar
	TokenMudar
	TokenPara

	// Relationship keywords
	TokenPertenceA
	TokenTemMuitos

	// Theme keywords
	TokenCor
	TokenIcone
	TokenEscuro

	// Modifier keywords
	TokenObrigatorio
	TokenUnico
	TokenPadrao
)

// Token represents a single lexical token.
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
	Indent  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%d, %q, line=%d, col=%d, indent=%d)", t.Type, t.Value, t.Line, t.Column, t.Indent)
}

var keywords = map[string]TokenType{
	// ==========================================
	// PORTUGUÊS
	// ==========================================

	// Blocos
	"sistema":     TokenSistema,
	"dados":       TokenDados,
	"telas":       TokenTelas,
	"acoes":       TokenAcoes,
	"eventos":     TokenEventos,
	"integracoes": TokenIntegracoes,
	"tema":        TokenTema,
	"logica":      TokenLogica,
	"banco":       TokenBanco,

	// Import
	"importar": TokenImportar,
	"de":       TokenDe,

	// Tela
	"tela":       TokenTela,
	"titulo":     TokenTitulo,
	"lista":      TokenLista,
	"mostrar":    TokenMostrar,
	"botao":      TokenBotao,
	"formulario": TokenFormulario,
	"entrada":    TokenEntrada,
	"busca":      TokenBusca,
	"dashboard":  TokenDashboard,

	// Tipos
	"texto":    TokenTexto,
	"numero":   TokenNumero,
	"data":     TokenData,
	"booleano": TokenBooleano,
	"email":    TokenEmail,
	"telefone": TokenTelefone,
	"imagem":   TokenImagem,
	"arquivo":  TokenArquivo,
	"upload":   TokenUpload,
	"link":     TokenLink,
	"status":   TokenStatus,
	"dinheiro": TokenDinheiro,
	"senha":    TokenSenha,

	// Eventos
	"quando":    TokenQuando,
	"clicar":    TokenClicar,
	"criar":     TokenCriar,
	"atualizar": TokenAtualizar,
	"deletar":   TokenDeletar,
	"enviar":    TokenEnviar,

	// Lógica
	"se":       TokenSe,
	"senao":    TokenSenao,
	"igual":    TokenIgual,
	"maior":    TokenMaior,
	"menor":    TokenMenor,
	"e":        TokenE,
	"ou":       TokenOu,
	"entao":    TokenEntao,
	"validar":  TokenValidar,
	"calcular": TokenCalcular,
	"definir":  TokenDefinir,
	"retornar": TokenRetornar,
	"mudar":    TokenMudar,
	"para":     TokenPara,

	// Relacionamentos
	"pertence_a": TokenPertenceA,
	"tem_muitos": TokenTemMuitos,

	// Tema
	"cor":    TokenCor,
	"icone":  TokenIcone,
	"escuro": TokenEscuro,

	// Modificadores
	"obrigatorio": TokenObrigatorio,
	"unico":       TokenUnico,
	"padrao":      TokenPadrao,

	// ==========================================
	// ENGLISH
	// ==========================================

	// Blocks
	"system":       TokenSistema,
	"models":       TokenDados,
	"screens":      TokenTelas,
	"actions":      TokenAcoes,
	"events":       TokenEventos,
	"integrations": TokenIntegracoes,
	"theme":        TokenTema,
	"logic":        TokenLogica,
	"database":     TokenBanco,
	"db":           TokenBanco,

	// Import
	"import": TokenImportar,
	"from":   TokenDe,

	// Screen
	"screen": TokenTela,
	"title":  TokenTitulo,
	"list":   TokenLista,
	"show":   TokenMostrar,
	"button": TokenBotao,
	"form":   TokenFormulario,
	"input":  TokenEntrada,
	"search": TokenBusca,

	// Types
	"text":     TokenTexto,
	"number":   TokenNumero,
	"date":     TokenData,
	"boolean":  TokenBooleano,
	"phone":    TokenTelefone,
	"image":    TokenImagem,
	"file":     TokenArquivo,
	"money":    TokenDinheiro,
	"password": TokenSenha,
	"currency": TokenDinheiro,

	// Events
	"when":   TokenQuando,
	"click":  TokenClicar,
	"create": TokenCriar,
	"update": TokenAtualizar,
	"delete": TokenDeletar,
	"send":   TokenEnviar,

	// Logic
	"if":       TokenSe,
	"else":     TokenSenao,
	"equals":   TokenIgual,
	"equal":    TokenIgual,
	"greater":  TokenMaior,
	"less":     TokenMenor,
	"and":      TokenE,
	"or":       TokenOu,
	"then":     TokenEntao,
	"validate": TokenValidar,
	"compute":  TokenCalcular,
	"set":      TokenDefinir,
	"return":   TokenRetornar,
	"change":   TokenMudar,
	"to":       TokenPara,

	// Relationships
	"belongs_to": TokenPertenceA,
	"has_many":   TokenTemMuitos,

	// Theme
	"color": TokenCor,
	"icon":  TokenIcone,
	"dark":  TokenEscuro,

	// Modifiers
	"required": TokenObrigatorio,
	"unique":   TokenUnico,
	"default":  TokenPadrao,

	// Color names (work in both languages)
	"azul":     TokenIdentifier,
	"verde":    TokenIdentifier,
	"vermelho": TokenIdentifier,
	"blue":     TokenIdentifier,
	"green":    TokenIdentifier,
	"red":      TokenIdentifier,
}

// Lexer tokenizes Flang source code.
type Lexer struct {
	source  []rune
	pos     int
	line    int
	col     int
	tokens  []Token
}

// New creates a new Lexer for the given source code.
func New(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
		pos:    0,
		line:   1,
		col:    1,
	}
}

// Tokenize processes the entire source and returns all tokens.
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.source) {
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}
	l.tokens = append(l.tokens, Token{Type: TokenEOF, Line: l.line, Column: l.col})
	return l.tokens, nil
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

func (l *Lexer) advance() rune {
	ch := l.source[l.pos]
	l.pos++
	l.col++
	return ch
}

func (l *Lexer) scanToken() error {
	ch := l.peek()

	// Handle newlines
	if ch == '\n' {
		l.tokens = append(l.tokens, Token{Type: TokenNewline, Value: "\n", Line: l.line, Column: l.col})
		l.advance()
		l.line++
		l.col = 1

		// Count indentation on the new line
		indent := 0
		for l.pos < len(l.source) && (l.peek() == ' ' || l.peek() == '\t') {
			if l.peek() == '\t' {
				indent += 2
			} else {
				indent++
			}
			l.advance()
		}
		if l.pos < len(l.source) && l.peek() != '\n' && l.peek() != '\r' {
			l.tokens = append(l.tokens, Token{Type: TokenIndent, Value: fmt.Sprintf("%d", indent), Line: l.line, Column: 1, Indent: indent})
		}
		return nil
	}

	// Skip carriage return
	if ch == '\r' {
		l.advance()
		return nil
	}

	// Skip spaces (mid-line)
	if ch == ' ' || ch == '\t' {
		l.advance()
		return nil
	}

	// Skip comments (// or #)
	if ch == '/' && l.pos+1 < len(l.source) && l.source[l.pos+1] == '/' {
		for l.pos < len(l.source) && l.peek() != '\n' {
			l.advance()
		}
		return nil
	}
	if ch == '#' {
		for l.pos < len(l.source) && l.peek() != '\n' {
			l.advance()
		}
		return nil
	}

	// Colon
	if ch == ':' {
		l.tokens = append(l.tokens, Token{Type: TokenColon, Value: ":", Line: l.line, Column: l.col})
		l.advance()
		return nil
	}

	// Dot
	if ch == '.' {
		l.tokens = append(l.tokens, Token{Type: TokenDot, Value: ".", Line: l.line, Column: l.col})
		l.advance()
		return nil
	}

	// String literal
	if ch == '"' {
		return l.scanString()
	}

	// Number
	if unicode.IsDigit(ch) {
		return l.scanNumber()
	}

	// Identifier or keyword
	if unicode.IsLetter(ch) || ch == '_' {
		return l.scanIdentifier()
	}

	return fmt.Errorf("unexpected character %q at line %d, column %d", string(ch), l.line, l.col)
}

func (l *Lexer) scanString() error {
	startLine := l.line
	startCol := l.col
	l.advance() // skip opening quote
	var buf strings.Builder
	for l.pos < len(l.source) && l.peek() != '"' {
		if l.peek() == '\n' {
			return fmt.Errorf("unterminated string at line %d, column %d", startLine, startCol)
		}
		if l.peek() == '\\' {
			l.advance()
			if l.pos >= len(l.source) {
				return fmt.Errorf("unterminated escape at line %d", l.line)
			}
			switch l.peek() {
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case '"':
				buf.WriteByte('"')
			case '\\':
				buf.WriteByte('\\')
			default:
				buf.WriteRune(l.peek())
			}
			l.advance()
			continue
		}
		buf.WriteRune(l.advance())
	}
	if l.pos >= len(l.source) {
		return fmt.Errorf("unterminated string at line %d, column %d", startLine, startCol)
	}
	l.advance() // skip closing quote
	l.tokens = append(l.tokens, Token{Type: TokenString, Value: buf.String(), Line: startLine, Column: startCol})
	return nil
}

func (l *Lexer) scanNumber() error {
	startCol := l.col
	var buf strings.Builder
	for l.pos < len(l.source) && (unicode.IsDigit(l.peek()) || l.peek() == '.') {
		buf.WriteRune(l.advance())
	}
	l.tokens = append(l.tokens, Token{Type: TokenNumber, Value: buf.String(), Line: l.line, Column: startCol})
	return nil
}

func (l *Lexer) scanIdentifier() error {
	startCol := l.col
	var buf strings.Builder
	for l.pos < len(l.source) && (unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_') {
		buf.WriteRune(l.advance())
	}
	word := buf.String()
	lower := strings.ToLower(word)

	if tt, ok := keywords[lower]; ok {
		l.tokens = append(l.tokens, Token{Type: tt, Value: lower, Line: l.line, Column: startCol})
	} else {
		l.tokens = append(l.tokens, Token{Type: TokenIdentifier, Value: word, Line: l.line, Column: startCol})
	}
	return nil
}

// IsBlockKeyword returns true if the token type is a top-level block keyword.
func IsBlockKeyword(tt TokenType) bool {
	switch tt {
	case TokenSistema, TokenDados, TokenTelas, TokenAcoes, TokenEventos, TokenIntegracoes, TokenTema, TokenLogica, TokenBanco:
		return true
	}
	return false
}

// IsTypeKeyword returns true if the token is a data type.
func IsTypeKeyword(tt TokenType) bool {
	switch tt {
	case TokenTexto, TokenNumero, TokenData, TokenBooleano, TokenEmail,
		TokenTelefone, TokenImagem, TokenArquivo, TokenUpload, TokenLink,
		TokenStatus, TokenDinheiro, TokenSenha:
		return true
	}
	return false
}
