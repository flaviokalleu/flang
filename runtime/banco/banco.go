package banco

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flavio/flang/compiler/ast"
	_ "modernc.org/sqlite"
)

// Banco wraps the database connection and model metadata.
type Banco struct {
	DB     *sql.DB
	Models map[string]*ast.Model
}

// Abrir creates the database and tables from model definitions.
func Abrir(path string, models []*ast.Model) (*Banco, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

	b := &Banco{
		DB:     db,
		Models: make(map[string]*ast.Model),
	}

	for _, m := range models {
		b.Models[strings.ToLower(m.Name)] = m
		if err := b.criarTabela(m); err != nil {
			return nil, fmt.Errorf("erro ao criar tabela '%s': %w", m.Name, err)
		}
		fmt.Printf("[flang] Tabela: %s (%d campos)\n", m.Name, len(m.Fields))
	}

	return b, nil
}

// q quotes a SQL identifier to avoid conflicts with reserved words.
func q(name string) string {
	return `"` + name + `"`
}

func (b *Banco) criarTabela(model *ast.Model) error {
	name := strings.ToLower(model.Name)

	var cols []string
	cols = append(cols, q("id")+" INTEGER PRIMARY KEY AUTOINCREMENT")
	for _, f := range model.Fields {
		col := fmt.Sprintf("%s %s", q(strings.ToLower(f.Name)), f.Type.SQLType())
		if f.Required {
			col += " NOT NULL"
		}
		if f.Unique {
			col += " UNIQUE"
		}
		if f.Default != "" {
			col += fmt.Sprintf(" DEFAULT '%s'", f.Default)
		}
		if f.Reference != "" {
			col += fmt.Sprintf(" REFERENCES %s(%s)", q(strings.ToLower(f.Reference)), q("id"))
		}
		cols = append(cols, col)
	}
	cols = append(cols, q("criado_em")+" DATETIME DEFAULT CURRENT_TIMESTAMP")
	cols = append(cols, q("atualizado_em")+" DATETIME DEFAULT CURRENT_TIMESTAMP")

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)", q(name), strings.Join(cols, ",\n  "))
	_, err := b.DB.Exec(query)
	return err
}

// Validar checks field constraints before insert/update.
func (b *Banco) Validar(modelo string, dados map[string]any) error {
	model, ok := b.Models[modelo]
	if !ok {
		return fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	for _, f := range model.Fields {
		fname := strings.ToLower(f.Name)
		val, exists := dados[fname]

		// Check required
		if f.Required && (!exists || val == nil || val == "") {
			return fmt.Errorf("campo '%s' é obrigatório", f.Name)
		}

		if !exists || val == nil {
			continue
		}

		strVal := fmt.Sprintf("%v", val)

		// Validate email format
		if f.Type == ast.FieldEmail && strVal != "" {
			if !strings.Contains(strVal, "@") || !strings.Contains(strVal, ".") {
				return fmt.Errorf("email inválido no campo '%s'", f.Name)
			}
		}

		// Validate phone format
		if f.Type == ast.FieldTelefone && strVal != "" {
			clean := strings.Map(func(r rune) rune {
				if r >= '0' && r <= '9' || r == '+' {
					return r
				}
				return -1
			}, strVal)
			if len(clean) < 7 {
				return fmt.Errorf("telefone inválido no campo '%s'", f.Name)
			}
		}
	}

	return nil
}

// Listar returns all rows from a model table.
func (b *Banco) Listar(modelo string) ([]map[string]any, error) {
	model, ok := b.Models[modelo]
	if !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	rows, err := b.DB.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC", q(modelo), q("id")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	results := []map[string]any{}

	for rows.Next() {
		values := make([]any, len(columns))
		ptrs := make([]any, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := map[string]any{}
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	_ = model
	return results, nil
}

// Buscar returns a single row by ID.
func (b *Banco) Buscar(modelo string, id int64) (map[string]any, error) {
	if _, ok := b.Models[modelo]; !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	rows, err := b.DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", q(modelo), q("id")), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	if !rows.Next() {
		return nil, fmt.Errorf("registro %d não encontrado", id)
	}

	values := make([]any, len(columns))
	ptrs := make([]any, len(columns))
	for i := range values {
		ptrs[i] = &values[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return nil, err
	}

	row := map[string]any{}
	for i, col := range columns {
		row[col] = values[i]
	}
	return row, nil
}

// Criar inserts a new row from JSON data.
func (b *Banco) Criar(modelo string, dados json.RawMessage) (map[string]any, error) {
	model, ok := b.Models[modelo]
	if !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	var input map[string]any
	if err := json.Unmarshal(dados, &input); err != nil {
		return nil, fmt.Errorf("dados inválidos: %w", err)
	}

	// Validate
	if err := b.Validar(modelo, input); err != nil {
		return nil, err
	}

	var cols []string
	var placeholders []string
	var vals []any

	for _, f := range model.Fields {
		fname := strings.ToLower(f.Name)
		if v, exists := input[fname]; exists {
			cols = append(cols, q(fname))
			placeholders = append(placeholders, "?")
			vals = append(vals, v)
		}
	}

	if len(cols) == 0 {
		return nil, fmt.Errorf("nenhum campo fornecido")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		q(modelo), strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	result, err := b.DB.Exec(query, vals...)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return b.Buscar(modelo, id)
}

// Atualizar updates a row by ID.
func (b *Banco) Atualizar(modelo string, id int64, dados json.RawMessage) (map[string]any, error) {
	model, ok := b.Models[modelo]
	if !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	var input map[string]any
	if err := json.Unmarshal(dados, &input); err != nil {
		return nil, fmt.Errorf("dados inválidos: %w", err)
	}

	var sets []string
	var vals []any

	for _, f := range model.Fields {
		fname := strings.ToLower(f.Name)
		if v, exists := input[fname]; exists {
			sets = append(sets, q(fname)+" = ?")
			vals = append(vals, v)
		}
	}

	if len(sets) == 0 {
		return nil, fmt.Errorf("nenhum campo para atualizar")
	}

	sets = append(sets, q("atualizado_em")+" = CURRENT_TIMESTAMP")
	vals = append(vals, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", q(modelo), strings.Join(sets, ", "), q("id"))
	_, err := b.DB.Exec(query, vals...)
	if err != nil {
		return nil, err
	}

	return b.Buscar(modelo, id)
}

// Deletar removes a row by ID.
func (b *Banco) Deletar(modelo string, id int64) error {
	if _, ok := b.Models[modelo]; !ok {
		return fmt.Errorf("modelo '%s' não encontrado", modelo)
	}
	_, err := b.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = ?", q(modelo), q("id")), id)
	return err
}

// Fechar closes the database connection.
func (b *Banco) Fechar() {
	b.DB.Close()
}
