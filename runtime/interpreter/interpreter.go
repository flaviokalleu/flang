package interpreter

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flavio/flang/compiler/ast"
	"github.com/flavio/flang/runtime/banco"
)

const maxIterations = 10000

// signalType is used for control flow via panic/recover.
type signalType int

const (
	signalReturn   signalType = iota
	signalBreak
	signalContinue
)

// signal carries a control flow signal and an optional value.
type signal struct {
	Type  signalType
	Value interface{}
}

// Scope holds variable bindings with parent scope chaining.
type Scope struct {
	vars   map[string]interface{}
	parent *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{vars: make(map[string]interface{}), parent: parent}
}

func (s *Scope) Get(name string) (interface{}, bool) {
	if v, ok := s.vars[name]; ok {
		return v, true
	}
	if s.parent != nil {
		return s.parent.Get(name)
	}
	return nil, false
}

func (s *Scope) Set(name string, value interface{}) {
	// Update in the scope where it exists, or set in current scope.
	if _, ok := s.vars[name]; ok {
		s.vars[name] = value
		return
	}
	if s.parent != nil {
		if _, ok := s.parent.Get(name); ok {
			s.parent.Set(name, value)
			return
		}
	}
	s.vars[name] = value
}

func (s *Scope) SetLocal(name string, value interface{}) {
	s.vars[name] = value
}

// Interpreter executes Flang AST scripts.
type Interpreter struct {
	Global    *Scope
	Functions map[string]*ast.FuncDecl
	DB        *banco.Banco
	LogBuffer []string
	logMu     sync.Mutex
}

// New creates a new interpreter.
func New(db *banco.Banco) *Interpreter {
	return &Interpreter{
		Global:    NewScope(nil),
		Functions: make(map[string]*ast.FuncDecl),
		DB:        db,
	}
}

// AppendLog adds a message to the log buffer.
func (interp *Interpreter) AppendLog(msg string) {
	interp.logMu.Lock()
	defer interp.logMu.Unlock()
	interp.LogBuffer = append(interp.LogBuffer, msg)
	// Keep only last 1000 log entries
	if len(interp.LogBuffer) > 1000 {
		interp.LogBuffer = interp.LogBuffer[len(interp.LogBuffer)-1000:]
	}
}

// GetLogs returns and optionally clears the log buffer.
func (interp *Interpreter) GetLogs(clear bool) []string {
	interp.logMu.Lock()
	defer interp.logMu.Unlock()
	logs := make([]string, len(interp.LogBuffer))
	copy(logs, interp.LogBuffer)
	if clear {
		interp.LogBuffer = nil
	}
	return logs
}

// RegisterFunction registers a function declaration.
func (interp *Interpreter) RegisterFunction(fn *ast.FuncDecl) {
	interp.Functions[fn.Name] = fn
}

// ExecStatements executes a list of statements in the given scope.
func (interp *Interpreter) ExecStatements(stmts []*ast.Statement, scope *Scope) {
	for _, stmt := range stmts {
		interp.ExecStatement(stmt, scope)
	}
}

// ExecStatement executes a single statement.
func (interp *Interpreter) ExecStatement(stmt *ast.Statement, scope *Scope) {
	if stmt == nil {
		return
	}

	switch stmt.Type {
	case "var":
		val := interp.EvalExpr(&stmt.VarDecl.Value, scope)
		scope.Set(stmt.VarDecl.Name, val)

	case "assign":
		val := interp.EvalExpr(&stmt.Assign.Value, scope)
		if stmt.Assign.Field != "" {
			// Object field assignment
			obj, ok := scope.Get(stmt.Assign.Target)
			if !ok {
				obj = make(map[string]interface{})
				scope.Set(stmt.Assign.Target, obj)
			}
			if m, ok := obj.(map[string]interface{}); ok {
				m[stmt.Assign.Field] = val
			}
		} else {
			scope.Set(stmt.Assign.Target, val)
		}

	case "if":
		interp.execIf(stmt.If, scope)

	case "for_each":
		interp.execForEach(stmt.ForEach, scope)

	case "while":
		interp.execWhile(stmt.While, scope)

	case "repeat":
		interp.execRepeat(stmt.Repeat, scope)

	case "return":
		val := interp.EvalExpr(stmt.Return, scope)
		panic(signal{Type: signalReturn, Value: val})

	case "break":
		panic(signal{Type: signalBreak})

	case "continue":
		panic(signal{Type: signalContinue})

	case "pause":
		time.Sleep(1 * time.Second)

	case "print":
		val := interp.EvalExpr(stmt.Print, scope)
		msg := toString(val)
		fmt.Println("[flang]", msg)
		interp.AppendLog(msg)

	case "call":
		interp.execCall(stmt.Call, scope)

	case "try":
		interp.execTry(stmt.Try, scope)
	}
}

func (interp *Interpreter) execIf(ifStmt *ast.IfStmt, scope *Scope) {
	if toBool(interp.EvalExpr(&ifStmt.Condition, scope)) {
		interp.ExecStatements(ifStmt.Body, NewScope(scope))
		return
	}

	for _, elseIf := range ifStmt.ElseIfs {
		if toBool(interp.EvalExpr(&elseIf.Condition, scope)) {
			interp.ExecStatements(elseIf.Body, NewScope(scope))
			return
		}
	}

	if len(ifStmt.Else) > 0 {
		interp.ExecStatements(ifStmt.Else, NewScope(scope))
	}
}

func (interp *Interpreter) execForEach(forEach *ast.ForEachStmt, scope *Scope) {
	collection := interp.EvalExpr(&forEach.Collection, scope)

	var items []interface{}
	switch c := collection.(type) {
	case []interface{}:
		items = c
	case []map[string]interface{}:
		for _, m := range c {
			items = append(items, m)
		}
	case string:
		// If it looks like a model name, query the DB
		if interp.DB != nil {
			if _, ok := interp.DB.Models[strings.ToLower(c)]; ok {
				rows, _, err := interp.DB.Listar(strings.ToLower(c), nil)
				if err == nil {
					for _, r := range rows {
						items = append(items, r)
					}
				}
			}
		}
		if items == nil {
			// Iterate over characters
			for _, ch := range c {
				items = append(items, string(ch))
			}
		}
	default:
		return
	}

	loopScope := NewScope(scope)
	iterations := 0
	for _, item := range items {
		iterations++
		if iterations > maxIterations {
			interp.AppendLog("ERRO: limite de iteracoes atingido (for_each)")
			break
		}
		loopScope.SetLocal(forEach.VarName, item)
		shouldBreak := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					if sig, ok := r.(signal); ok {
						switch sig.Type {
						case signalBreak:
							shouldBreak = true
							return
						case signalContinue:
							return
						default:
							panic(r) // re-panic return signals
						}
					}
					panic(r)
				}
			}()
			interp.ExecStatements(forEach.Body, loopScope)
		}()
		if shouldBreak {
			break
		}
	}
}

func (interp *Interpreter) execWhile(whileStmt *ast.WhileStmt, scope *Scope) {
	iterations := 0
	loopScope := NewScope(scope)
	for {
		iterations++
		if iterations > maxIterations {
			interp.AppendLog("ERRO: limite de iteracoes atingido (while)")
			break
		}
		if !toBool(interp.EvalExpr(&whileStmt.Condition, loopScope)) {
			break
		}
		shouldBreak := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					if sig, ok := r.(signal); ok {
						switch sig.Type {
						case signalBreak:
							shouldBreak = true
							return
						case signalContinue:
							return
						default:
							panic(r)
						}
					}
					panic(r)
				}
			}()
			interp.ExecStatements(whileStmt.Body, loopScope)
		}()
		if shouldBreak {
			break
		}
	}
}

func (interp *Interpreter) execRepeat(repeatStmt *ast.RepeatStmt, scope *Scope) {
	countVal := interp.EvalExpr(&repeatStmt.Count, scope)
	count := int(toNumber(countVal))
	if count > maxIterations {
		count = maxIterations
		interp.AppendLog("AVISO: repeat limitado a 10000 iteracoes")
	}

	loopScope := NewScope(scope)
	for i := 0; i < count; i++ {
		loopScope.SetLocal("i", float64(i))
		shouldBreak := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					if sig, ok := r.(signal); ok {
						switch sig.Type {
						case signalBreak:
							shouldBreak = true
							return
						case signalContinue:
							return
						default:
							panic(r)
						}
					}
					panic(r)
				}
			}()
			interp.ExecStatements(repeatStmt.Body, loopScope)
		}()
		if shouldBreak {
			break
		}
	}
}

func (interp *Interpreter) execCall(call *ast.FuncCall, scope *Scope) interface{} {
	if call == nil {
		return nil
	}

	name := call.Name
	args := make([]interface{}, len(call.Args))
	for i, arg := range call.Args {
		args[i] = interp.EvalExpr(arg, scope)
	}

	// Check built-in functions
	if result, ok := interp.callBuiltin(name, args); ok {
		return result
	}

	// Check user-defined functions
	if fn, ok := interp.Functions[name]; ok {
		return interp.callFunction(fn, args)
	}

	// DB operations via object.method pattern
	if call.Object != "" {
		return interp.callDBMethod(call.Object, name, args, scope)
	}

	return nil
}

func (interp *Interpreter) callFunction(fn *ast.FuncDecl, args []interface{}) interface{} {
	fnScope := NewScope(interp.Global)
	for i, param := range fn.Params {
		if i < len(args) {
			fnScope.SetLocal(param, args[i])
		} else {
			fnScope.SetLocal(param, nil)
		}
	}

	var result interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if sig, ok := r.(signal); ok && sig.Type == signalReturn {
					result = sig.Value
					return
				}
				panic(r)
			}
		}()
		interp.ExecStatements(fn.Body, fnScope)
	}()

	return result
}

func (interp *Interpreter) callDBMethod(object, method string, args []interface{}, scope *Scope) interface{} {
	if interp.DB == nil {
		return nil
	}

	modelName := strings.ToLower(object)
	if _, ok := interp.DB.Models[modelName]; !ok {
		return nil
	}

	switch method {
	case "listar", "list", "todos", "all":
		rows, _, err := interp.DB.Listar(modelName, nil)
		if err != nil {
			interp.AppendLog(fmt.Sprintf("ERRO DB listar: %s", err))
			return nil
		}
		result := make([]interface{}, len(rows))
		for i, r := range rows {
			result[i] = r
		}
		return result

	case "criar", "create":
		if len(args) > 0 {
			data, err := json.Marshal(args[0])
			if err != nil {
				interp.AppendLog(fmt.Sprintf("ERRO DB criar: %s", err))
				return nil
			}
			row, err := interp.DB.Criar(modelName, data)
			if err != nil {
				interp.AppendLog(fmt.Sprintf("ERRO DB criar: %s", err))
				return nil
			}
			return row
		}

	case "atualizar", "update":
		if len(args) >= 2 {
			id := int64(toNumber(args[0]))
			data, err := json.Marshal(args[1])
			if err != nil {
				return nil
			}
			row, err := interp.DB.Atualizar(modelName, id, data)
			if err != nil {
				interp.AppendLog(fmt.Sprintf("ERRO DB atualizar: %s", err))
				return nil
			}
			return row
		}

	case "deletar", "delete":
		if len(args) > 0 {
			id := int64(toNumber(args[0]))
			if err := interp.DB.Deletar(modelName, id); err != nil {
				interp.AppendLog(fmt.Sprintf("ERRO DB deletar: %s", err))
				return nil
			}
			return true
		}

	case "contar", "count":
		rows, total, err := interp.DB.Listar(modelName, nil)
		if err != nil {
			return float64(0)
		}
		if total > 0 {
			return float64(total)
		}
		return float64(len(rows))
	}

	return nil
}

func (interp *Interpreter) execTry(tryStmt *ast.TryStmt, scope *Scope) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if it's a control flow signal — re-panic those
				if sig, ok := r.(signal); ok {
					panic(sig)
				}
				// Handle the error
				if len(tryStmt.Catch) > 0 {
					catchScope := NewScope(scope)
					if tryStmt.ErrVar != "" {
						catchScope.SetLocal(tryStmt.ErrVar, fmt.Sprintf("%v", r))
					}
					interp.ExecStatements(tryStmt.Catch, catchScope)
				}
			}
		}()
		interp.ExecStatements(tryStmt.Body, NewScope(scope))
	}()
}

// EvalExpr evaluates an expression and returns the result.
func (interp *Interpreter) EvalExpr(expr *ast.Expression, scope *Scope) interface{} {
	if expr == nil {
		return nil
	}

	switch expr.Type {
	case "literal":
		// Parse number strings to float64
		if s, ok := expr.Value.(string); ok {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return f
			}
		}
		return expr.Value

	case "variable":
		if val, ok := scope.Get(expr.Name); ok {
			return val
		}
		// Could be a model name — return as string
		return expr.Name

	case "binary":
		return interp.evalBinary(expr, scope)

	case "unary":
		return interp.evalUnary(expr, scope)

	case "call":
		call := &ast.FuncCall{
			Name:   expr.Name,
			Object: expr.Object,
		}
		args := make([]interface{}, len(expr.Args))
		for i, a := range expr.Args {
			args[i] = interp.EvalExpr(a, scope)
		}
		// Built-in functions
		if result, ok := interp.callBuiltin(expr.Name, args); ok {
			return result
		}
		// User-defined functions
		if fn, ok := interp.Functions[expr.Name]; ok {
			return interp.callFunction(fn, args)
		}
		// Object method calls
		if expr.Object != "" {
			call.Args = expr.Args
			result := interp.callDBMethod(expr.Object, expr.Name, args, scope)
			if result != nil {
				return result
			}
			// Try as a method on a variable
			if obj, ok := scope.Get(expr.Object); ok {
				return interp.callObjMethod(obj, expr.Name, args)
			}
		}
		return nil

	case "field_access":
		if obj, ok := scope.Get(expr.Object); ok {
			if m, ok := obj.(map[string]interface{}); ok {
				return m[expr.Field]
			}
		}
		// Try as model name for DB access
		if interp.DB != nil {
			modelName := strings.ToLower(expr.Object)
			if _, ok := interp.DB.Models[modelName]; ok {
				// Return a string reference for lazy evaluation
				return fmt.Sprintf("%s.%s", expr.Object, expr.Field)
			}
		}
		return nil

	case "list":
		result := make([]interface{}, len(expr.Elements))
		for i, elem := range expr.Elements {
			result[i] = interp.EvalExpr(elem, scope)
		}
		return result
	}

	return nil
}

func (interp *Interpreter) evalBinary(expr *ast.Expression, scope *Scope) interface{} {
	left := interp.EvalExpr(expr.Left, scope)
	right := interp.EvalExpr(expr.Right, scope)

	switch expr.Operator {
	case "+":
		// String concatenation if either is string
		if isString(left) || isString(right) {
			return toString(left) + toString(right)
		}
		return toNumber(left) + toNumber(right)
	case "-":
		return toNumber(left) - toNumber(right)
	case "*":
		return toNumber(left) * toNumber(right)
	case "/":
		r := toNumber(right)
		if r == 0 {
			return float64(0)
		}
		return toNumber(left) / r
	case "==":
		return isEqual(left, right)
	case "!=":
		return !isEqual(left, right)
	case ">":
		return toNumber(left) > toNumber(right)
	case "<":
		return toNumber(left) < toNumber(right)
	case ">=":
		return toNumber(left) >= toNumber(right)
	case "<=":
		return toNumber(left) <= toNumber(right)
	case "e":
		return toBool(left) && toBool(right)
	case "ou":
		return toBool(left) || toBool(right)
	}

	return nil
}

func (interp *Interpreter) evalUnary(expr *ast.Expression, scope *Scope) interface{} {
	val := interp.EvalExpr(expr.Right, scope)
	switch expr.Operator {
	case "nao":
		return !toBool(val)
	case "-":
		return -toNumber(val)
	}
	return nil
}

func (interp *Interpreter) callObjMethod(obj interface{}, method string, args []interface{}) interface{} {
	switch method {
	case "tamanho", "length", "len":
		switch v := obj.(type) {
		case string:
			return float64(len(v))
		case []interface{}:
			return float64(len(v))
		case map[string]interface{}:
			return float64(len(v))
		}
	}
	return nil
}

// ==================== Built-in Functions ====================

func (interp *Interpreter) callBuiltin(name string, args []interface{}) (interface{}, bool) {
	switch name {
	case "tamanho", "length", "len":
		if len(args) < 1 {
			return float64(0), true
		}
		switch v := args[0].(type) {
		case string:
			return float64(len(v)), true
		case []interface{}:
			return float64(len(v)), true
		case map[string]interface{}:
			return float64(len(v)), true
		default:
			return float64(0), true
		}

	case "tipo", "type":
		if len(args) < 1 {
			return "nulo", true
		}
		switch args[0].(type) {
		case float64:
			return "numero", true
		case string:
			return "texto", true
		case bool:
			return "booleano", true
		case nil:
			return "nulo", true
		case []interface{}:
			return "lista", true
		case map[string]interface{}:
			return "objeto", true
		default:
			return "desconhecido", true
		}

	case "texto", "string":
		if len(args) < 1 {
			return "", true
		}
		return toString(args[0]), true

	case "numero", "number":
		if len(args) < 1 {
			return float64(0), true
		}
		return toNumber(args[0]), true

	case "arredondar", "round":
		if len(args) < 1 {
			return float64(0), true
		}
		return math.Round(toNumber(args[0])), true

	case "aleatorio", "random":
		return rand.Float64(), true

	case "agora", "now":
		return time.Now().Format(time.RFC3339), true

	case "maiusculo", "uppercase", "upper":
		if len(args) < 1 {
			return "", true
		}
		return strings.ToUpper(toString(args[0])), true

	case "minusculo", "lowercase", "lower":
		if len(args) < 1 {
			return "", true
		}
		return strings.ToLower(toString(args[0])), true

	case "contem", "contains":
		if len(args) < 2 {
			return false, true
		}
		return strings.Contains(toString(args[0]), toString(args[1])), true

	case "dividir", "split":
		if len(args) < 2 {
			return []interface{}{}, true
		}
		parts := strings.Split(toString(args[0]), toString(args[1]))
		result := make([]interface{}, len(parts))
		for i, p := range parts {
			result[i] = p
		}
		return result, true

	case "juntar", "join":
		if len(args) < 2 {
			return "", true
		}
		if arr, ok := args[0].([]interface{}); ok {
			strs := make([]string, len(arr))
			for i, v := range arr {
				strs[i] = toString(v)
			}
			return strings.Join(strs, toString(args[1])), true
		}
		return "", true

	case "abs":
		if len(args) < 1 {
			return float64(0), true
		}
		return math.Abs(toNumber(args[0])), true

	case "min":
		if len(args) < 2 {
			return float64(0), true
		}
		return math.Min(toNumber(args[0]), toNumber(args[1])), true

	case "max":
		if len(args) < 2 {
			return float64(0), true
		}
		return math.Max(toNumber(args[0]), toNumber(args[1])), true

	case "inteiro", "int":
		if len(args) < 1 {
			return float64(0), true
		}
		return math.Floor(toNumber(args[0])), true
	}

	return nil, false
}

// ==================== Type Helpers ====================

func toNumber(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return f
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return "nulo"
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == math.Floor(val) && !math.IsInf(val, 0) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		if val {
			return "verdadeiro"
		}
		return "falso"
	case []interface{}:
		parts := make([]string, len(val))
		for i, v := range val {
			parts[i] = toString(v)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]interface{}:
		data, _ := json.Marshal(val)
		return string(data)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	case string:
		return val != "" && val != "0" && val != "falso" && val != "false"
	case []interface{}:
		return len(val) > 0
	case map[string]interface{}:
		return len(val) > 0
	default:
		return true
	}
}

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Try numeric comparison
	aNum, aIsNum := tryNumber(a)
	bNum, bIsNum := tryNumber(b)
	if aIsNum && bIsNum {
		return aNum == bNum
	}
	// Fall back to string comparison
	return toString(a) == toString(b)
}

func tryNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

// Run executes an entire program's scripts and registers functions.
func (interp *Interpreter) Run(program *ast.Program) {
	// Register functions
	for _, fn := range program.Functions {
		interp.RegisterFunction(fn)
	}

	// Execute top-level scripts
	func() {
		defer func() {
			if r := recover(); r != nil {
				if sig, ok := r.(signal); ok && sig.Type == signalReturn {
					return // top-level return is fine
				}
				interp.AppendLog(fmt.Sprintf("ERRO runtime: %v", r))
				fmt.Printf("[flang] ERRO runtime: %v\n", r)
			}
		}()
		interp.ExecStatements(program.Scripts, interp.Global)
	}()
}

// Eval parses and executes a code snippet, returning output.
func (interp *Interpreter) Eval(code string) ([]string, error) {
	// We need to import the lexer and parser here
	// But to avoid circular imports, we accept pre-parsed statements
	return nil, fmt.Errorf("use EvalStatements instead")
}

// EvalStatements executes statements and returns the log output.
func (interp *Interpreter) EvalStatements(stmts []*ast.Statement, fns []*ast.FuncDecl) []string {
	// Register any new functions
	for _, fn := range fns {
		interp.RegisterFunction(fn)
	}

	// Clear logs for this eval
	interp.logMu.Lock()
	startIdx := len(interp.LogBuffer)
	interp.logMu.Unlock()

	func() {
		defer func() {
			if r := recover(); r != nil {
				if sig, ok := r.(signal); ok && sig.Type == signalReturn {
					return
				}
				interp.AppendLog(fmt.Sprintf("ERRO: %v", r))
			}
		}()
		interp.ExecStatements(stmts, interp.Global)
	}()

	interp.logMu.Lock()
	defer interp.logMu.Unlock()
	if startIdx >= len(interp.LogBuffer) {
		return nil
	}
	return interp.LogBuffer[startIdx:]
}
