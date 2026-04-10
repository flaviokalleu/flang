package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flavio/flang/compiler/ast"
	"github.com/flavio/flang/compiler/lexer"
	"github.com/flavio/flang/compiler/parser"
	"github.com/flavio/flang/runtime/banco"
	"github.com/flavio/flang/runtime/servidor"
	wa "github.com/flavio/flang/runtime/whatsapp"
)

// parseFG reads and parses a single .fg file.
func parseFG(arquivo string) (*ast.Program, error) {
	source, err := os.ReadFile(arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler %s: %w", arquivo, err)
	}

	lex := lexer.New(string(source))
	tokens, err := lex.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("erro léxico em %s: %w", arquivo, err)
	}

	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("erro de parsing em %s: %w", arquivo, err)
	}

	return program, nil
}

// resolveImports processes all import statements recursively.
func resolveImports(program *ast.Program, baseDir string, resolved map[string]bool) error {
	if resolved == nil {
		resolved = make(map[string]bool)
	}

	for _, imp := range program.Imports {
		// Resolve path relative to the base .fg file
		importPath := filepath.Join(baseDir, imp.Path)
		absPath, err := filepath.Abs(importPath)
		if err != nil {
			return fmt.Errorf("caminho inválido: %s", imp.Path)
		}

		// Avoid circular imports
		if resolved[absPath] {
			continue
		}
		resolved[absPath] = true

		fmt.Printf("[flang] Importando: %s\n", imp.Path)

		imported, err := parseFG(absPath)
		if err != nil {
			return fmt.Errorf("erro ao importar %s: %w", imp.Path, err)
		}

		// Recursively resolve imports in the imported file
		importDir := filepath.Dir(absPath)
		if err := resolveImports(imported, importDir, resolved); err != nil {
			return err
		}

		// Merge based on what was requested
		switch imp.What {
		case "tudo":
			program.Merge(imported)
		case "dados":
			program.Models = append(program.Models, imported.Models...)
		case "telas":
			program.Screens = append(program.Screens, imported.Screens...)
		case "eventos":
			program.Events = append(program.Events, imported.Events...)
		case "tema":
			if imported.Theme != nil {
				program.Theme = imported.Theme
			}
		case "logica":
			program.Rules = append(program.Rules, imported.Rules...)
		default:
			// Import specific named items (e.g., importar produto de "dados.fg")
			for _, m := range imported.Models {
				if m.Name == imp.What {
					program.Models = append(program.Models, m)
				}
			}
			for _, s := range imported.Screens {
				if s.Name == imp.What {
					program.Screens = append(program.Screens, s)
				}
			}
		}
	}

	return nil
}

// Executar loads a .fg file and runs the application.
func Executar(arquivo string, porta string) error {
	fmt.Printf("[flang] Carregando: %s\n", arquivo)

	program, err := parseFG(arquivo)
	if err != nil {
		return err
	}

	// Resolve imports
	baseDir := filepath.Dir(arquivo)
	if err := resolveImports(program, baseDir, nil); err != nil {
		return err
	}

	if program.System == nil {
		return fmt.Errorf("declaração 'sistema' não encontrada")
	}

	fmt.Printf("[flang] Sistema: %s\n", program.System.Name)
	fmt.Printf("[flang] Modelos: %d | Telas: %d | Eventos: %d | Regras: %d\n",
		len(program.Models), len(program.Screens), len(program.Events), len(program.Rules))

	// Database
	db, err := banco.Abrir(program.Database, program.System.Name, program.Models)
	if err != nil {
		return fmt.Errorf("erro no banco: %w", err)
	}

	// WhatsApp
	var waClient *wa.Client
	if program.WhatsApp != nil && program.WhatsApp.Enabled {
		waClient = wa.Novo(program.WhatsApp.DBPath)
		if err := waClient.Conectar(); err != nil {
			fmt.Printf("[flang] AVISO WhatsApp: %s (continuando sem WhatsApp)\n", err)
			waClient = nil
		}
		if waClient != nil {
			defer waClient.Desconectar()
		}
	}

	// Server
	srv := servidor.Novo(program, db, porta)
	srv.WA = waClient
	fmt.Printf("\n[flang] %s rodando em http://localhost:%s\n\n", program.System.Name, porta)
	return srv.Iniciar()
}

// Verificar loads and parses a .fg file without running.
func Verificar(arquivo string) error {
	program, err := parseFG(arquivo)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(arquivo)
	if err := resolveImports(program, baseDir, nil); err != nil {
		return err
	}

	if program.System == nil {
		return fmt.Errorf("declaração 'sistema' não encontrada")
	}

	fmt.Printf("[flang] ✓ %s - válido\n", arquivo)
	fmt.Printf("  sistema:  %s\n", program.System.Name)
	fmt.Printf("  imports:  %d\n", len(program.Imports))
	fmt.Printf("  modelos:  %d\n", len(program.Models))
	fmt.Printf("  telas:    %d\n", len(program.Screens))
	fmt.Printf("  eventos:  %d\n", len(program.Events))
	fmt.Printf("  regras:   %d\n", len(program.Rules))
	if program.WhatsApp != nil && program.WhatsApp.Enabled {
		fmt.Printf("  whatsapp: ativado (%d notificações)\n", len(program.Notifiers))
	}
	return nil
}
