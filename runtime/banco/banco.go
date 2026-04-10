package banco

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flavio/flang/compiler/ast"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// Banco wraps the database connection and model metadata.
type Banco struct {
	DB     *sql.DB
	Models map[string]*ast.Model
	Driver string // "sqlite", "mysql", "postgres"
}

// Abrir creates the database and tables from model definitions.
func Abrir(config *ast.DatabaseConfig, appName string, models []*ast.Model) (*Banco, error) {
	if config == nil {
		config = ast.DefaultDatabase()
	}

	driver, dsn := buildDSN(config, appName)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar (%s): %w", config.Driver, err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco %s: %w", config.Driver, err)
	}

	fmt.Printf("[flang] Banco: %s\n", config.Driver)

	b := &Banco{
		DB:     db,
		Models: make(map[string]*ast.Model),
		Driver: config.Driver,
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

func buildDSN(config *ast.DatabaseConfig, appName string) (driver string, dsn string) {
	switch strings.ToLower(config.Driver) {
	case "mysql":
		host := config.Host
		if host == "" {
			host = "localhost"
		}
		port := config.Port
		if port == "" {
			port = "3306"
		}
		dbName := config.Name
		if dbName == "" {
			dbName = appName
		}
		user := config.User
		if user == "" {
			user = "root"
		}
		// user:password@tcp(host:port)/dbname?parseTime=true
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			user, config.Password, host, port, dbName)
		return "mysql", dsn

	case "postgres", "postgresql":
		host := config.Host
		if host == "" {
			host = "localhost"
		}
		port := config.Port
		if port == "" {
			port = "5432"
		}
		dbName := config.Name
		if dbName == "" {
			dbName = appName
		}
		user := config.User
		if user == "" {
			user = "postgres"
		}
		sslmode := "disable"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, config.Password, dbName, sslmode)
		return "postgres", dsn

	default: // sqlite
		dbName := config.Name
		if dbName == "" {
			dbName = appName + ".db"
		}
		return "sqlite", dbName + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)"
	}
}

// q quotes a SQL identifier.
func q(name string) string {
	return `"` + name + `"`
}

// placeholder returns the correct placeholder for the driver.
func (b *Banco) ph(n int) string {
	if b.Driver == "postgres" || b.Driver == "postgresql" {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
}

// placeholders returns N placeholders for the driver.
func (b *Banco) placeholders(count int) []string {
	phs := make([]string, count)
	for i := range phs {
		phs[i] = b.ph(i + 1)
	}
	return phs
}

func (b *Banco) criarTabela(model *ast.Model) error {
	name := strings.ToLower(model.Name)

	autoInc := "AUTOINCREMENT"
	if b.Driver == "mysql" {
		autoInc = "AUTO_INCREMENT"
	}
	if b.Driver == "postgres" || b.Driver == "postgresql" {
		autoInc = "" // use SERIAL instead
	}

	var cols []string
	if b.Driver == "postgres" || b.Driver == "postgresql" {
		cols = append(cols, q("id")+" SERIAL PRIMARY KEY")
	} else {
		cols = append(cols, q("id")+" INTEGER PRIMARY KEY "+autoInc)
	}

	for _, f := range model.Fields {
		sqlType := f.Type.SQLType()
		// Adjust types per driver
		if b.Driver == "mysql" {
			if sqlType == "TEXT" {
				sqlType = "VARCHAR(500)"
			}
		}

		col := fmt.Sprintf("%s %s", q(strings.ToLower(f.Name)), sqlType)
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

	tsDefault := "CURRENT_TIMESTAMP"
	tsType := "DATETIME"
	if b.Driver == "postgres" || b.Driver == "postgresql" {
		tsType = "TIMESTAMP"
	}

	cols = append(cols, fmt.Sprintf("%s %s DEFAULT %s", q("criado_em"), tsType, tsDefault))
	cols = append(cols, fmt.Sprintf("%s %s DEFAULT %s", q("atualizado_em"), tsType, tsDefault))

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)", q(name), strings.Join(cols, ",\n  "))
	_, err := b.DB.Exec(query)
	return err
}

// Listar returns all rows from a model table.
func (b *Banco) Listar(modelo string) ([]map[string]any, error) {
	if _, ok := b.Models[modelo]; !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	rows, err := b.DB.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC", q(modelo), q("id")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

// Buscar returns a single row by ID.
func (b *Banco) Buscar(modelo string, id int64) (map[string]any, error) {
	if _, ok := b.Models[modelo]; !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	rows, err := b.DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", q(modelo), q("id"), b.ph(1)), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := scanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("registro %d não encontrado", id)
	}
	return results[0], nil
}

// Criar inserts a new row.
func (b *Banco) Criar(modelo string, dados json.RawMessage) (map[string]any, error) {
	model, ok := b.Models[modelo]
	if !ok {
		return nil, fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	var input map[string]any
	if err := json.Unmarshal(dados, &input); err != nil {
		return nil, fmt.Errorf("dados inválidos: %w", err)
	}

	if err := b.Validar(modelo, input); err != nil {
		return nil, err
	}

	var cols []string
	var phs []string
	var vals []any
	n := 1

	for _, f := range model.Fields {
		fname := strings.ToLower(f.Name)
		if v, exists := input[fname]; exists {
			cols = append(cols, q(fname))
			phs = append(phs, b.ph(n))
			vals = append(vals, v)
			n++
		}
	}

	if len(cols) == 0 {
		return nil, fmt.Errorf("nenhum campo fornecido")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		q(modelo), strings.Join(cols, ", "), strings.Join(phs, ", "))

	if b.Driver == "postgres" || b.Driver == "postgresql" {
		query += " RETURNING " + q("id")
		var id int64
		err := b.DB.QueryRow(query, vals...).Scan(&id)
		if err != nil {
			return nil, err
		}
		return b.Buscar(modelo, id)
	}

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
	n := 1

	for _, f := range model.Fields {
		fname := strings.ToLower(f.Name)
		if v, exists := input[fname]; exists {
			sets = append(sets, q(fname)+" = "+b.ph(n))
			vals = append(vals, v)
			n++
		}
	}

	if len(sets) == 0 {
		return nil, fmt.Errorf("nenhum campo para atualizar")
	}

	sets = append(sets, q("atualizado_em")+" = CURRENT_TIMESTAMP")
	vals = append(vals, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s",
		q(modelo), strings.Join(sets, ", "), q("id"), b.ph(n))
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
	_, err := b.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = %s", q(modelo), q("id"), b.ph(1)), id)
	return err
}

// Fechar closes the database connection.
func (b *Banco) Fechar() {
	b.DB.Close()
}

// Contar returns the count of rows in a table.
func (b *Banco) Contar(modelo string) (int64, error) {
	var count int64
	err := b.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", q(modelo))).Scan(&count)
	return count, err
}

// Validar checks field constraints.
func (b *Banco) Validar(modelo string, dados map[string]any) error {
	model, ok := b.Models[modelo]
	if !ok {
		return fmt.Errorf("modelo '%s' não encontrado", modelo)
	}

	for _, f := range model.Fields {
		fname := strings.ToLower(f.Name)
		val, exists := dados[fname]

		if f.Required && (!exists || val == nil || val == "") {
			return fmt.Errorf("campo '%s' é obrigatório", f.Name)
		}

		if !exists || val == nil {
			continue
		}

		strVal := fmt.Sprintf("%v", val)

		if f.Type == ast.FieldEmail && strVal != "" {
			if !strings.Contains(strVal, "@") || !strings.Contains(strVal, ".") {
				return fmt.Errorf("email inválido no campo '%s'", f.Name)
			}
		}

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

// scanRows converts sql.Rows into a slice of maps.
func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	columns, _ := rows.Columns()
	var results []map[string]any

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

	return results, nil
}
