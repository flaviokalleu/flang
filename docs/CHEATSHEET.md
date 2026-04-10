# Flang - Cheatsheet / Referencia Rapida

## Estrutura Minima

```
sistema nome
dados
  modelo
    campo: tipo
telas
  tela nome
    titulo "Titulo"
    lista modelo
      mostrar campo
    botao azul
      texto "Acao"
eventos
  quando clicar "Acao"
    criar modelo
```

## Tipos

```
texto / text          numero / number        dinheiro / money
email                 telefone / phone       status
data / date           booleano / boolean     senha / password
imagem / image        arquivo / file         upload
link
```

## Modificadores

```
campo: tipo obrigatorio / required
campo: tipo unico / unique
campo: tipo pertence_a modelo / belongs_to model
```

## Telas

```
tela nome / screen name
  titulo "X" / title "X"
  lista modelo / list model
    mostrar campo / show field
  botao azul / button blue
    texto "X" / text "X"
```

## Eventos

```
quando clicar "X" / when click "X"
  criar modelo / create model
```

## Tema

```
tema / theme
  cor primaria "#hex" / color primary "#hex"
  cor secundaria "#hex" / color secondary "#hex"
  cor destaque "#hex" / color accent "#hex"
  escuro / dark
```

## Banco

```
banco / database
  driver: sqlite / mysql / postgres
  host: "localhost"
  porta / port: "5432"
  nome / name: "db"
  usuario / user: "user"
  senha / password: "pass"
```

## Imports

```
importar "arquivo.fg" / import "file.fg"
importar dados de "x.fg" / import models from "x.fg"
```

## Logica

```
logica / logic
  validar campo condicao / validate field condition
  se campo igual "valor" / if field equals "value"
    mudar acao / change action
```

## WhatsApp

```
integracoes / integrations
  whatsapp
    quando criar modelo / when create model
      enviar mensagem para campo / send message to field
        texto "Msg {var}" / text "Msg {var}"
```

## CLI

```bash
flang run arquivo.fg [porta]
flang check arquivo.fg
flang new nome
flang version
flang help
```

## API REST (automatica)

```
GET    /api/{modelo}      Listar
GET    /api/{modelo}/{id}  Buscar
POST   /api/{modelo}      Criar
PUT    /api/{modelo}/{id}  Atualizar
DELETE /api/{modelo}/{id}  Deletar
```

## WebSocket

```
ws://localhost:8080/ws
```
