# Script para popular um banco de dados com filmes

## Tecnologias usadas

- Golang v1.20.0
- Docker Compose v2.15.1
- PostgreSQL latest

## Como rodar esse projeto de forma local

1. Clonar esse repositório pelo terminal
   - Via HTTP
     `git clone https://github.com/RickChaves29/script_populate_db.git`

   - Via SSH
     `git clone git@github.com:RickChaves29/script_populate_db.git`

2. Ainda no terminal, copie a variável de ambiente que está no arquivo .env.example e cole no arquivo .bashrc ou .profile adicionando a palavra chave **export** antes.

   > OBS: O Arquivo **.bashrc** fica na pasta raiz do seu úsuario

   - Exemplo no WSL ou Linux:

     ```bash
     export CONNECT_DB='postgres://user:password@host:port/dbname?sslmode=disable'
     ```

3. Voltando para pasta onde você clonou o projeto rode os seguintes comandos:

    >OBS: A tabela no banco de dados irar ser criada automaticamente pelo script
   - Crie um aquivo chamado .env na pasta do projeto e copie duas varieveis de ambiente, que são:

   ```bash
    POSTGRES_USER='user'
    POSTGRES_PASSWORD='password'
   ```

   - Para subir o banco de dados `docker compose up`
   - Rodar o script `go run script.go`
