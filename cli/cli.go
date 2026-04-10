package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"

	"github.com/flavio/flang/runtime"
)

const version = "0.2.0"

const banner = `
  ███████╗██╗      █████╗ ███╗   ██╗ ██████╗
  ██╔════╝██║     ██╔══██╗████╗  ██║██╔════╝
  █████╗  ██║     ███████║██╔██╗ ██║██║  ███╗
  ██╔══╝  ██║     ██╔══██╗██║╚██╗██║██║   ██║
  ██║     ███████╗██║  ██║██║ ╚████║╚██████╔╝
  ╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝
  v%s - Tudo roda direto do .fg
`

// Run executes the CLI.
func Run(args []string) {
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch args[1] {
	case "run":
		if len(args) < 3 {
			fmt.Println("Uso: flang run <arquivo.fg>")
			os.Exit(1)
		}
		porta := "8080"
		if len(args) >= 4 {
			porta = args[3]
		}
		if err := runtime.Executar(args[2], porta); err != nil {
			fmt.Printf("[flang] ERRO: %s\n", err)
			os.Exit(1)
		}

	case "check":
		if len(args) < 3 {
			fmt.Println("Uso: flang check <arquivo.fg>")
			os.Exit(1)
		}
		if err := runtime.Verificar(args[2]); err != nil {
			fmt.Printf("[flang] ERRO: %s\n", err)
			os.Exit(1)
		}

	case "new":
		if len(args) < 3 {
			fmt.Println("Uso: flang new <nome>")
			os.Exit(1)
		}
		cmdNew(args[2])

	case "version":
		fmt.Printf(banner, version)

	case "docker":
		cmdDocker()

	case "init":
		if len(args) < 3 {
			fmt.Println("Uso: flang init <nome>")
			os.Exit(1)
		}
		cmdInit(args[2])

	case "build":
		if len(args) < 3 {
			fmt.Println("Uso: flang build <arquivo.fg> [--output nome]")
			os.Exit(1)
		}
		output := ""
		for i, a := range args {
			if (a == "--output" || a == "-o") && i+1 < len(args) {
				output = args[i+1]
			}
		}
		cmdBuild(args[2], output)

	case "help":
		printUsage()

	default:
		// If arg ends in .fg, treat it as "run"
		if strings.HasSuffix(args[1], ".fg") {
			porta := "8080"
			if len(args) >= 3 {
				porta = args[2]
			}
			if err := runtime.Executar(args[1], porta); err != nil {
				fmt.Printf("[flang] ERRO: %s\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Comando desconhecido: %s\n", args[1])
			printUsage()
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Printf(banner, version)
	fmt.Println(`
Uso: flang <comando> [argumentos]

Comandos:
  run <arquivo.fg> [porta]  Executa o arquivo .fg (porta padrao: 8080)
  check <arquivo.fg>        Verifica sintaxe sem executar
  new <nome>                Cria projeto plano (tudo num arquivo so)
  init <nome>               Cria projeto organizado (pastas por responsabilidade)
  build <arquivo.fg> [-o nome]  Compila em executavel standalone
  docker                    Gera Dockerfile para o projeto atual
  version                   Mostra a versao
  help                      Mostra esta ajuda

Modos de projeto:
  new   → Modo plano: um arquivo so, ideal para comecar rapido.
  init  → Modo organizado: dados/, telas/, eventos/ separados.
          Comece com 'new' e migre para 'init' quando crescer.

Atalho:
  flang inicio.fg           Mesmo que "flang run inicio.fg"

Exemplo:
  flang new meuapp          Cria projeto plano
  flang init meuapp         Cria projeto organizado
  flang run meuapp/inicio.fg
`)
}

func cmdNew(name string) {
	title := strings.ToUpper(name[:1]) + name[1:]
	dir := name
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Erro: %s\n", err)
		os.Exit(1)
	}

	// Modo plano: tudo num arquivo só, simples e direto
	fg := `sistema ` + name + `

tema
  cor primaria "#6366f1"
  cor secundaria "#8b5cf6"

dados

  produto
    nome: texto obrigatorio
    descricao: texto
    preco: dinheiro
    estoque: numero
    status: status

  cliente
    nome: texto obrigatorio
    email: email unico
    telefone: telefone
    status: status

telas

  tela produtos
    titulo "Produtos"
    lista produto
      mostrar nome
      mostrar preco
      mostrar estoque
      mostrar status
    botao azul
      texto "Novo Produto"

  tela clientes
    titulo "Clientes"
    lista cliente
      mostrar nome
      mostrar email
      mostrar telefone
      mostrar status
    botao verde
      texto "Novo Cliente"

eventos

  quando clicar "Novo Produto"
    criar produto

  quando clicar "Novo Cliente"
    criar cliente
`

	fgPath := filepath.Join(dir, "inicio.fg")
	if err := os.WriteFile(fgPath, []byte(fg), 0644); err != nil {
		fmt.Printf("Erro: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("[flang] Projeto '%s' criado! (modo plano)\n", title)
	fmt.Println("[flang] Tudo num arquivo so - simples e direto.")
	fmt.Printf("[flang] Execute: flang %s\n", fgPath)
	fmt.Println()
	fmt.Println("[flang] Dica: quando crescer, use 'flang init' para modo organizado.")
}

func cmdDocker() {
	// Find .fg files in the current directory
	fgFile := "inicio.fg"
	entries, err := os.ReadDir(".")
	if err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".fg") {
				fgFile = e.Name()
				break
			}
		}
	}

	dockerfile := fmt.Sprintf(`# Generated by flang docker
FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o flang .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /build/flang /usr/local/bin/flang
COPY *.fg ./

EXPOSE 8080
CMD ["flang", "run", "%s"]
`, fgFile)

	if err := os.WriteFile("Dockerfile", []byte(dockerfile), 0644); err != nil {
		fmt.Printf("[flang] Erro ao criar Dockerfile: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("[flang] Dockerfile gerado com sucesso!")
	fmt.Println("[flang] Execute: docker build -t meu-app . && docker run -p 8080:8080 meu-app")
}

func cmdInit(name string) {
	title := strings.ToUpper(name[:1]) + name[1:]
	dir := name

	// Criar estrutura organizada por responsabilidade
	// Inspirado no React: cada pasta tem um papel claro
	dirs := []string{
		dir,
		filepath.Join(dir, "dados"),  // modelos (como models/ ou types/)
		filepath.Join(dir, "telas"),  // interfaces (como pages/ ou components/)
		filepath.Join(dir, "eventos"), // interacoes (como handlers/ ou hooks/)
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			fmt.Printf("Erro: %s\n", err)
			os.Exit(1)
		}
	}

	// ── inicio.fg ── entry point (como App.js no React)
	inicio := `sistema ` + name + `

importar "tema.fg"
importar "dados/produto.fg"
importar "dados/cliente.fg"
importar "telas/produtos.fg"
importar "telas/clientes.fg"
importar "eventos/acoes.fg"
`
	// ── tema.fg ── visual (como theme.js)
	tema := `tema
  cor primaria "#6366f1"
  cor secundaria "#8b5cf6"
  cor destaque "#f59e0b"
`

	// ── dados/produto.fg ── um modelo por arquivo (como um component)
	produto := `dados

  produto
    nome: texto obrigatorio
    descricao: texto
    preco: dinheiro
    estoque: numero
    categoria: texto
    status: status
`

	// ── dados/cliente.fg
	cliente := `dados

  cliente
    nome: texto obrigatorio
    email: email unico
    telefone: telefone
    cidade: texto
    status: status
`

	// ── telas/produtos.fg
	telaProdutos := `telas

  tela produtos
    titulo "Produtos"
    lista produto
      mostrar nome
      mostrar preco
      mostrar estoque
      mostrar categoria
      mostrar status
    botao azul
      texto "Novo Produto"
`

	// ── telas/clientes.fg
	telaClientes := `telas

  tela clientes
    titulo "Clientes"
    lista cliente
      mostrar nome
      mostrar email
      mostrar telefone
      mostrar cidade
      mostrar status
    botao verde
      texto "Novo Cliente"
`

	// ── eventos/acoes.fg
	acoes := `eventos

  quando clicar "Novo Produto"
    criar produto

  quando clicar "Novo Cliente"
    criar cliente
`

	// Mapa de arquivos a criar
	files := map[string]string{
		filepath.Join(dir, "inicio.fg"):           inicio,
		filepath.Join(dir, "tema.fg"):              tema,
		filepath.Join(dir, "dados", "produto.fg"):  produto,
		filepath.Join(dir, "dados", "cliente.fg"):  cliente,
		filepath.Join(dir, "telas", "produtos.fg"): telaProdutos,
		filepath.Join(dir, "telas", "clientes.fg"): telaClientes,
		filepath.Join(dir, "eventos", "acoes.fg"):  acoes,
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("Erro ao criar %s: %s\n", path, err)
			os.Exit(1)
		}
	}

	// ── .env
	envContent := `# Configuracao do projeto ` + name + `
FLANG_PORT=8080
FLANG_DB_TYPE=sqlite
FLANG_DB_NAME=` + name + `.db
`
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		fmt.Printf("Erro ao criar .env: %s\n", err)
		os.Exit(1)
	}

	// ── .gitignore
	gitignore := `*.db
*.db-shm
*.db-wal
.env
flang
flang.exe
`
	giPath := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(giPath, []byte(gitignore), 0644); err != nil {
		fmt.Printf("Erro ao criar .gitignore: %s\n", err)
		os.Exit(1)
	}

	// ── Dockerfile
	dockerfileContent := `FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o flang .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /build/flang /usr/local/bin/flang
COPY *.fg ./
COPY dados/ ./dados/
COPY telas/ ./telas/
COPY eventos/ ./eventos/

EXPOSE 8080
CMD ["flang", "run", "inicio.fg"]
`
	dfPath := filepath.Join(dir, "Dockerfile")
	if err := os.WriteFile(dfPath, []byte(dockerfileContent), 0644); err != nil {
		fmt.Printf("Erro ao criar Dockerfile: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("[flang] Projeto '%s' criado! (modo organizado)\n", title)
	fmt.Println()
	fmt.Printf("  %s/\n", name)
	fmt.Printf("  ├── inicio.fg          (entry point)\n")
	fmt.Printf("  ├── tema.fg            (visual)\n")
	fmt.Printf("  ├── dados/\n")
	fmt.Printf("  │   ├── produto.fg     (modelo)\n")
	fmt.Printf("  │   └── cliente.fg     (modelo)\n")
	fmt.Printf("  ├── telas/\n")
	fmt.Printf("  │   ├── produtos.fg    (interface)\n")
	fmt.Printf("  │   └── clientes.fg    (interface)\n")
	fmt.Printf("  ├── eventos/\n")
	fmt.Printf("  │   └── acoes.fg       (interacoes)\n")
	fmt.Printf("  ├── .env\n")
	fmt.Printf("  ├── .gitignore\n")
	fmt.Printf("  └── Dockerfile\n")
	fmt.Println()
	fmt.Printf("[flang] Execute: flang run %s\n", filepath.Join(name, "inicio.fg"))
	fmt.Println()
	fmt.Println("[flang] Adicione novos modelos em dados/, telas em telas/,")
	fmt.Println("        e importe no inicio.fg. Cada arquivo cuida de uma coisa.")
}

func cmdBuild(arquivo string, output string) {
	// Verify the .fg file exists
	if _, err := os.Stat(arquivo); os.IsNotExist(err) {
		fmt.Printf("[flang] Erro: arquivo '%s' nao encontrado\n", arquivo)
		os.Exit(1)
	}

	// First, verify the .fg file is valid
	if err := runtime.Verificar(arquivo); err != nil {
		fmt.Printf("[flang] Erro: %s\n", err)
		os.Exit(1)
	}

	// Determine output name
	if output == "" {
		base := filepath.Base(arquivo)
		output = strings.TrimSuffix(base, filepath.Ext(base))
		if goruntime.GOOS == "windows" {
			output += ".exe"
		}
	}

	// Collect all .fg files in the directory
	dir := filepath.Dir(arquivo)
	if dir == "" || dir == "." {
		dir, _ = os.Getwd()
	} else {
		dir, _ = filepath.Abs(dir)
	}
	mainFile := filepath.Base(arquivo)

	// Create temp build directory
	tmpDir, err := os.MkdirTemp("", "flang-build-*")
	if err != nil {
		fmt.Printf("[flang] Erro ao criar diretorio temporario: %s\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Copy all .fg files to temp dir preserving structure
	fgFiles := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".fg" {
			rel, _ := filepath.Rel(dir, path)
			destPath := filepath.Join(tmpDir, "app", rel)
			os.MkdirAll(filepath.Dir(destPath), 0755)
			data, _ := os.ReadFile(path)
			os.WriteFile(destPath, data, 0644)
			fgFiles = append(fgFiles, rel)
		}
		return nil
	})

	// Copy .env if exists
	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		data, _ := os.ReadFile(envPath)
		os.MkdirAll(filepath.Join(tmpDir, "app"), 0755)
		os.WriteFile(filepath.Join(tmpDir, "app", ".env"), data, 0644)
	}

	// Generate main.go for the standalone binary
	mainGo := `package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flavio/flang/runtime"
)

//go:embed app/*
var appFS embed.FS

func main() {
	// Extract embedded files to temp dir
	tmpDir, err := os.MkdirTemp("", "flang-app-*")
	if err != nil {
		fmt.Printf("Erro: %s\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Walk embedded FS and write files
	entries, _ := appFS.ReadDir("app")
	extractDir(appFS, "app", tmpDir, entries)

	porta := "8080"
	if len(os.Args) > 1 {
		porta = os.Args[1]
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		porta = envPort
	}

	arquivo := filepath.Join(tmpDir, "` + mainFile + `")
	if err := runtime.Executar(arquivo, porta); err != nil {
		fmt.Printf("Erro: %s\n", err)
		os.Exit(1)
	}
}

func extractDir(fs embed.FS, base string, dest string, entries []os.DirEntry) {
	for _, e := range entries {
		srcPath := base + "/" + e.Name()
		destPath := filepath.Join(dest, e.Name())
		if e.IsDir() {
			os.MkdirAll(destPath, 0755)
			subEntries, _ := fs.ReadDir(srcPath)
			extractDir(fs, srcPath, destPath, subEntries)
		} else {
			data, _ := fs.ReadFile(srcPath)
			os.WriteFile(destPath, data, 0644)
		}
	}
}
`

	// Write main.go
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainGo), 0644); err != nil {
		fmt.Printf("[flang] Erro ao gerar main.go: %s\n", err)
		os.Exit(1)
	}

	// Find the flang module path for the replace directive
	flangModPath := getFlangModPath()

	// Generate go.mod
	goMod := "module flang-app\n\ngo 1.26\n\nrequire github.com/flavio/flang v0.0.0\n\nreplace github.com/flavio/flang => " + flangModPath + "\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		fmt.Printf("[flang] Erro ao gerar go.mod: %s\n", err)
		os.Exit(1)
	}

	// Build
	fmt.Printf("[flang] Compilando %s...\n", output)

	absOutput, _ := filepath.Abs(output)

	cmd := exec.Command("go", "build", "-o", absOutput, "-ldflags", "-s -w", ".")
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("[flang] Erro na compilacao: %s\n", err)
		os.Exit(1)
	}

	// Get file size
	info, _ := os.Stat(absOutput)
	sizeMB := float64(info.Size()) / (1024 * 1024)

	fmt.Printf("[flang] Build concluido: %s (%.1f MB)\n", output, sizeMB)
	fmt.Printf("[flang] Execute: ./%s [porta]\n", output)
}

func getFlangModPath() string {
	// Find the flang module path by looking for go.mod
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)

	// Walk up looking for go.mod with flang module
	for {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			data, _ := os.ReadFile(modPath)
			if strings.Contains(string(data), "flavio/flang") || strings.Contains(string(data), "module") {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback: try GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home, _ := os.UserHomeDir()
		gopath = filepath.Join(home, "go")
	}
	return filepath.Join(gopath, "src", "github.com", "flavio", "flang")
}
