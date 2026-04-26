# fc-go-expert-client-server

Desafio Go Expert — Client/Server com contextos, SQLite e timeouts.

## Pré-requisitos

- Go 1.21+
- (Opcional) Docker e Docker Compose

---

## Rodando localmente (sem Docker)

### 1. Instalar dependências

```bash
go mod tidy
```

### 2. Iniciar o servidor

```bash
go run server.go
```

O servidor sobe na porta **8080** e cria o arquivo `cotacoes.db` no diretório corrente.

### 3. Executar o cliente (em outro terminal)

```bash
go run client.go
```

O cliente salva o câmbio atual no arquivo **`cotacao.txt`** no formato `Dólar: {valor}`.

---

## Rodando com Docker Compose

```bash
docker-compose up --build
```

O banco de dados SQLite fica persistido no volume Docker `cotacoes-data` — os dados sobrevivem mesmo que o container seja removido.

Para executar o cliente apontando para o container:

```bash
go run client.go
```

---

## Depuração no VS Code

O arquivo `.vscode/launch.json` contém duas configurações:

| Configuração    | Descrição                              |
|-----------------|----------------------------------------|
| **Debug Server** | Inicia o `server.go` em modo debug    |
| **Debug Client** | Inicia o `client.go` em modo debug    |

Abra a aba **Run & Debug** (`Ctrl+Shift+D`), selecione a configuração desejada e pressione **F5**.

---

## Timeouts aplicados

| Operação                          | Timeout |
|-----------------------------------|---------|
| Chamada à AwesomeAPI (server)     | 200 ms  |
| Persistência no SQLite (server)   | 10 ms   |
| Resposta do servidor (client)     | 300 ms  |
