# Tutorial Completo - Flang

Aprenda a criar aplicacoes completas com Flang em menos de 30 minutos.

---

## Indice

1. [Instalacao](#1-instalacao)
2. [Seu Primeiro App](#2-seu-primeiro-app)
3. [Entendendo os Blocos](#3-entendendo-os-blocos)
4. [Modelos de Dados](#4-modelos-de-dados)
5. [Telas e Componentes](#5-telas-e-componentes)
6. [Eventos](#6-eventos)
7. [Tema e Visual](#7-tema-e-visual)
8. [Logica e Validacoes](#8-logica-e-validacoes)
9. [Sistema de Imports](#9-sistema-de-imports)
10. [Banco de Dados](#10-banco-de-dados)
11. [WebSocket](#11-websocket-tempo-real)
12. [WhatsApp](#12-whatsapp)
13. [API REST](#13-api-rest)
14. [English Mode](#14-english-mode)
15. [Projeto Completo](#15-projeto-completo)
16. [Deploy](#16-deploy)

---

## 1. Instalacao

### Requisitos

- Go 1.21 ou superior ([download](https://go.dev/dl/))

### Passo a passo

```bash
# Clone o repositorio
git clone https://github.com/flaviokalleu/flang.git

# Entre na pasta
cd flang

# Compile
go build -o flang .

# Teste
./flang version
```

Voce vera:
```
  ███████╗██╗      █████╗ ███╗   ██╗ ██████╗
  ██╔════╝██║     ██╔══██╗████╗  ██║██╔════╝
  █████╗  ██║     ███████║██╔██╗ ██║██║  ███╗
  ██╔══╝  ██║     ██╔══██╗██║╚██╗██║██║   ██║
  ██║     ███████╗██║  ██║██║ ╚████║╚██████╔╝
  ╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝
  v0.4 - Tudo roda direto do .fg
```

---

## 2. Seu Primeiro App

Crie um arquivo chamado `inicio.fg`:

```
sistema meu_app

dados

  tarefa
    titulo: texto obrigatorio
    descricao: texto
    status: status

telas

  tela tarefas

    titulo "Minhas Tarefas"

    lista tarefa

      mostrar titulo
      mostrar status

    botao azul
      texto "Nova Tarefa"

eventos

  quando clicar "Nova Tarefa"
    criar tarefa
```

Rode:

```bash
./flang run inicio.fg
```

Abra `http://localhost:8080` no navegador.

**Pronto!** Voce tem uma app de tarefas com:
- Dashboard com contagem
- Tabela com listagem
- Formulario para criar/editar
- Botoes de editar e excluir
- Badge de status colorido
- Dark mode
- Busca

---

## 3. Entendendo os Blocos

Um arquivo `.fg` e composto por **blocos**. Cada bloco define uma parte da aplicacao:

```
sistema nome_do_app       # Obrigatorio - define o nome

dados                      # Modelos de dados (tabelas)
  ...

telas                      # Interface do usuario
  ...

eventos                    # Acoes do usuario
  ...

tema                       # Cores e visual (opcional)
  ...

logica                     # Regras de negocio (opcional)
  ...

banco                      # Configuracao do banco (opcional)
  ...

integracoes                # WhatsApp, etc (opcional)
  ...
```

### Regras de sintaxe

- **Indentacao** define a hierarquia (use 2 espacos)
- **Sem ponto e virgula**, sem chaves, sem parenteses
- **Sem virgulas** entre campos
- **Strings** entre aspas duplas: `"texto aqui"`
- **Comentarios** com `#` ou `//`

---

## 4. Modelos de Dados

O bloco `dados` define suas tabelas. Cada modelo vira uma tabela no banco e uma API REST automatica.

### Sintaxe basica

```
dados

  nome_do_modelo
    campo: tipo
    campo: tipo modificador
```

### Todos os tipos

```
dados

  exemplo
    # Textos
    nome: texto                    # texto simples
    bio: texto obrigatorio         # texto obrigatorio

    # Numeros
    idade: numero                  # numero decimal
    preco: dinheiro                # formata como R$ 29.90

    # Contato
    email: email unico             # validacao de email + unique
    fone: telefone                 # validacao de telefone

    # Status
    situacao: status               # badge colorido automatico

    # Data
    nascimento: data               # input de data

    # Booleano
    ativo: booleano                # checkbox

    # Senha
    senha: senha                   # input mascarado

    # Arquivos
    foto: imagem                   # upload de imagem
    documento: arquivo             # upload de arquivo
    anexo: upload                  # upload generico

    # URL
    site: link                     # input de URL
```

### Modificadores

```
dados

  usuario
    nome: texto obrigatorio        # NOT NULL + validacao
    cpf: texto unico               # UNIQUE constraint
    empresa_id: numero pertence_a empresa  # Foreign Key
```

### Relacionamentos

```
dados

  empresa
    nome: texto obrigatorio

  funcionario
    nome: texto obrigatorio
    cargo: texto
    empresa_id: numero pertence_a empresa
```

O campo `empresa_id` cria uma foreign key para a tabela `empresa`.

---

## 5. Telas e Componentes

O bloco `telas` define a interface do usuario.

### Estrutura

```
telas

  tela nome_da_tela

    titulo "Titulo exibido"

    lista nome_do_modelo

      mostrar campo1
      mostrar campo2
      mostrar campo3

    botao cor
      texto "Texto do Botao"
```

### Cores de botao

```
botao azul          # azul/blue
botao verde         # verde/green
botao vermelho      # vermelho/red
```

### Multiplas telas

Cada `tela` vira uma secao na sidebar:

```
telas

  tela produtos
    titulo "Produtos"
    lista produto
      mostrar nome
      mostrar preco
    botao azul
      texto "Novo"

  tela clientes
    titulo "Clientes"
    lista cliente
      mostrar nome
      mostrar email
    botao verde
      texto "Novo"

  tela pedidos
    titulo "Pedidos"
    lista pedido
      mostrar produto
      mostrar status
    botao azul
      texto "Novo"
```

---

## 6. Eventos

O bloco `eventos` conecta acoes do usuario a operacoes:

```
eventos

  quando clicar "Novo Produto"
    criar produto

  quando clicar "Novo Cliente"
    criar cliente
```

O texto do `quando clicar` deve ser **exatamente** o mesmo do `texto` do botao.

---

## 7. Tema e Visual

Customize cores e visual:

```
tema

  cor primaria "#6366f1"       # cor principal (botoes, links)
  cor secundaria "#8b5cf6"     # cor de gradiente
  cor destaque "#f59e0b"       # cor de destaque (cards)
  cor sidebar "#1e1b4b"        # cor da sidebar
```

### Cores populares

| Estilo | Primaria | Secundaria | Destaque |
|--------|----------|------------|----------|
| Indigo (padrao) | `#6366f1` | `#8b5cf6` | `#f59e0b` |
| Azul | `#3b82f6` | `#2563eb` | `#f59e0b` |
| Verde | `#10b981` | `#059669` | `#6366f1` |
| Vermelho | `#ef4444` | `#dc2626` | `#f59e0b` |
| Rosa | `#ec4899` | `#db2777` | `#8b5cf6` |

---

## 8. Logica e Validacoes

O bloco `logica` define regras de negocio:

```
logica

  # Validacoes
  validar email obrigatorio unico
  validar preco maior 0

  # Condicoes
  se status igual "cancelado"
    mudar cor vermelho

  se status igual "ativo"
    mudar cor verde

  se quantidade maior 10
    validar observacao obrigatorio
```

---

## 9. Sistema de Imports

Divida projetos grandes em multiplos arquivos:

### Estrutura de pastas

```
meu-projeto/
  inicio.fg          # arquivo principal
  dados.fg           # modelos
  telas.fg           # telas
  eventos.fg         # eventos
  tema.fg            # visual
  regras.fg          # logica
```

### inicio.fg

```
sistema meu_projeto

importar "tema.fg"
importar "dados.fg"
importar "telas.fg"
importar "eventos.fg"
importar "regras.fg"
```

### Tipos de import

```
importar "arquivo.fg"              # importa tudo
importar dados de "modelos.fg"     # so os modelos
importar telas de "paginas.fg"     # so as telas
importar produto de "dados.fg"     # um modelo especifico
```

---

## 10. Banco de Dados

### SQLite (padrao)

Nao precisa configurar nada. O banco e criado automaticamente:

```
sistema meu_app
# Cria meu_app.db automaticamente
```

### PostgreSQL

```
banco
  driver: postgres
  host: "localhost"
  porta: "5432"
  nome: "meu_banco"
  usuario: "postgres"
  senha: "minhasenha"
```

### MySQL

```
banco
  driver: mysql
  host: "localhost"
  porta: "3306"
  nome: "meu_banco"
  usuario: "root"
  senha: "minhasenha"
```

### Atalho

```
banco postgres
  host: "localhost"
  nome: "meu_banco"
```

---

## 11. WebSocket (Tempo Real)

WebSocket esta **ativo automaticamente**. Nao precisa configurar nada.

### O que acontece

1. Usuario A cria um produto
2. Usuario B (em outra aba/computador) ve o produto aparecer instantaneamente
3. Sem refresh, sem polling

### Para desenvolvedores

O endpoint WebSocket e `ws://localhost:8080/ws`.

Mensagens recebidas:
```json
{"type": "criar", "model": "produto", "id": 1, "data": {...}}
{"type": "atualizar", "model": "produto", "id": 1, "data": {...}}
{"type": "deletar", "model": "produto", "id": 1}
```

---

## 12. WhatsApp

### Configuracao

Adicione o bloco `integracoes` com `whatsapp`:

```
integracoes

  whatsapp

    quando criar pedido
      enviar mensagem para telefone
        texto "Ola {cliente}! Pedido de {prato} recebido!"

    quando atualizar pedido
      enviar mensagem para telefone
        texto "Status do pedido: {status}"
```

### Primeiro uso

1. Rode `flang run inicio.fg`
2. Um QR Code aparece no terminal
3. No celular: WhatsApp > Dispositivos Conectados > Conectar Dispositivo
4. Escaneie o QR Code
5. Pronto - mensagens automaticas ativadas

### Templates

Use `{campo}` para inserir valores do registro:

```
texto "Ola {nome}! Seu pedido de {prato} no valor de R${valor} foi {status}."
```

### Destino

O campo `para` define para qual numero enviar:

```
enviar mensagem para telefone           # campo 'telefone' do registro
enviar mensagem para cliente.telefone   # campo de um modelo relacionado
```

---

## 13. API REST

Cada modelo gera uma API REST completa automaticamente.

### Endpoints

Para um modelo chamado `produto`:

| Metodo | Rota | Acao |
|--------|------|------|
| GET | `/api/produto` | Listar todos |
| GET | `/api/produto/1` | Buscar por ID |
| POST | `/api/produto` | Criar novo |
| PUT | `/api/produto/1` | Atualizar |
| DELETE | `/api/produto/1` | Deletar |

### Exemplos com curl

```bash
# Listar
curl http://localhost:8080/api/produto

# Criar
curl -X POST http://localhost:8080/api/produto \
  -H "Content-Type: application/json" \
  -d '{"nome":"Camiseta","preco":59.90,"status":"ativo"}'

# Atualizar
curl -X PUT http://localhost:8080/api/produto/1 \
  -H "Content-Type: application/json" \
  -d '{"preco":49.90,"status":"promocao"}'

# Deletar
curl -X DELETE http://localhost:8080/api/produto/1
```

### Validacoes na API

A API valida automaticamente:
- Campos `obrigatorio` retornam erro se vazios
- Campos `email` validam formato
- Campos `telefone` validam comprimento minimo
- Campos `unico` checam duplicatas

Exemplo de erro:
```json
{"erro": "campo 'nome' e obrigatorio"}
```

---

## 14. English Mode

Flang e totalmente bilingue. Todas as keywords tem equivalente em ingles:

```
system my_store

theme
  color primary "#3b82f6"
  dark

models
  product
    name: text required
    price: money
    status: status

screens
  screen products
    title "Products"
    list product
      show name
      show price
      show status
    button blue
      text "New Product"

events
  when click "New Product"
    create product

logic
  validate price greater 0
```

### Tabela completa PT → EN

| Portugues | English |
|-----------|---------|
| sistema | system |
| dados | models |
| telas | screens |
| eventos | events |
| tema | theme |
| logica | logic |
| banco | database |
| integracoes | integrations |
| tela | screen |
| titulo | title |
| lista | list |
| mostrar | show |
| botao | button |
| texto | text |
| numero | number |
| dinheiro | money |
| telefone | phone |
| imagem | image |
| arquivo | file |
| senha | password |
| quando | when |
| clicar | click |
| criar | create |
| atualizar | update |
| deletar | delete |
| enviar | send |
| se | if |
| senao | else |
| igual | equals |
| maior | greater |
| menor | less |
| validar | validate |
| mudar | change |
| obrigatorio | required |
| unico | unique |
| pertence_a | belongs_to |
| importar | import |
| de | from |
| cor | color |
| escuro | dark |
| mensagem | message |

### Misturando idiomas

Voce pode misturar livremente:

```
system pizzaria

dados
  pizza
    nome: text required
    preco: money

screens
  tela cardapio
    title "Cardapio"
    list pizza
      mostrar nome
      show preco
```

---

## 15. Projeto Completo

Aqui esta um projeto real de restaurante com todas as features:

### inicio.fg

```
sistema restaurante

importar "tema.fg"
importar "dados.fg"
importar "telas.fg"
importar "eventos.fg"

integracoes

  whatsapp

    quando criar pedido
      enviar mensagem para telefone
        texto "Pedido recebido! {prato} x{quantidade}"
```

### tema.fg

```
tema
  cor primaria "#6366f1"
  cor secundaria "#8b5cf6"
  cor destaque "#f59e0b"
```

### dados.fg

```
dados

  prato
    nome: texto obrigatorio
    preco: dinheiro obrigatorio
    categoria: texto
    status: status

  mesa
    numero: numero obrigatorio unico
    capacidade: numero
    status: status

  pedido
    mesa: numero
    prato: texto obrigatorio
    quantidade: numero obrigatorio
    telefone: telefone
    status: status

  funcionario
    nome: texto obrigatorio
    cargo: texto
    email: email unico
    telefone: telefone
    status: status
```

### telas.fg

```
telas

  tela cardapio
    titulo "Cardapio"
    lista prato
      mostrar nome
      mostrar preco
      mostrar categoria
      mostrar status
    botao azul
      texto "Novo Prato"

  tela mesas
    titulo "Mesas"
    lista mesa
      mostrar numero
      mostrar capacidade
      mostrar status
    botao verde
      texto "Nova Mesa"

  tela pedidos
    titulo "Pedidos"
    lista pedido
      mostrar prato
      mostrar quantidade
      mostrar status
    botao azul
      texto "Novo Pedido"

  tela equipe
    titulo "Equipe"
    lista funcionario
      mostrar nome
      mostrar cargo
      mostrar email
      mostrar status
    botao azul
      texto "Novo Funcionario"
```

### eventos.fg

```
eventos

  quando clicar "Novo Prato"
    criar prato

  quando clicar "Nova Mesa"
    criar mesa

  quando clicar "Novo Pedido"
    criar pedido

  quando clicar "Novo Funcionario"
    criar funcionario
```

---

## 16. Deploy

### Opcao 1: Binario direto

```bash
# Compile o Flang
go build -o flang .

# Copie flang + seus .fg para o servidor
scp flang inicio.fg user@server:~/app/

# No servidor
cd ~/app
./flang run inicio.fg 80
```

### Opcao 2: Com PostgreSQL em producao

```
# inicio.fg
sistema meu_app

banco
  driver: postgres
  host: "db.meuservidor.com"
  porta: "5432"
  nome: "producao"
  usuario: "app_user"
  senha: "senha_segura"

importar "dados.fg"
importar "telas.fg"
importar "eventos.fg"
```

### Opcao 3: Porta customizada

```bash
./flang run inicio.fg 3000
```

---

## Proximo Passos

- Explore os exemplos em `exemplos/`
- Crie seu proprio projeto com `./flang new meu-projeto`
- Contribua no [GitHub](https://github.com/flaviokalleu/flang)

---

<p align="center">
  <strong>Flang</strong> - Descreva. Execute. Pronto.
</p>
