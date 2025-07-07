# Process Worker - Hackathon SOAT

Este projeto implementa o `hackthon-soat-process-worker`, worker responsável por processar vídeos de forma assíncrona para entrega do hackathon da pós graduação de Arquitetura de Software da FIAP.

O serviço opera em segundo plano, consumindo tarefas de uma fila de mensagens. Para cada tarefa, ele baixa um vídeo de um serviço de armazenamento, extrai todos os frames usando FFmpeg, compacta os frames em um único arquivo `.zip` e disponibiliza o resultado final para download.

### ✨ Arquitetura

O sistema foi projetado seguindo as melhores práticas de engenharia de software:
* **Microserviços:** O worker é um componente independente, desacoplado de outras partes do sistema (como a API).
* **Arquitetura Hexagonal:** A lógica de negócio (o "core") é completamente isolada de detalhes de infraestrutura (banco de dados, filas, etc.), promovendo alta testabilidade e manutenibilidade.
* **Orientado a Mensageria:** A comunicação é feita de forma assíncrona através de filas de mensagens (SQS), o que torna o sistema resiliente a picos de carga.

### 🛠️ Tecnologias Utilizadas

* **Linguagem:** Go (1.23)
* **Banco de Dados:** PostgreSQL
* **ORM:** GORM
* **Cache:** Redis
* **Fila de Mensagens:** Amazon SQS (simulado com LocalStack)
* **Armazenamento de Arquivos:** Amazon S3 (simulado com LocalStack)
* **Processamento de Vídeo:** FFmpeg
* **Automação de Tarefas:** Makefile

---

## 🚀 Execução do Ambiente Local

Siga estes passos para iniciar o ambiente de desenvolvimento completo na sua máquina.

### 1. Pré-requisitos

Antes de começar, garanta que você tem as seguintes ferramentas instaladas:

* **Git:** Para clonar o projeto.
* **Go:** Versão 1.23 ou superior.
* **Docker & Docker Desktop:** Essencial para rodar os contêineres da nossa infraestrutura.
* **Make:** Para executar os comandos de atalho do `Makefile`.
    * **Windows:** A forma mais fácil é usar o terminal **Git Bash**, que já vem com o `make`. Alternativamente, instale via Chocolatey (`choco install make`).
* **AWS CLI:** Necessário para validar os recursos no LocalStack manualmente.

### 2. Configuração Inicial (Apenas na primeira vez)

1.  **Clone o repositório** para a sua máquina.

2.  **Crie o arquivo de ambiente local:** Navegue até a pasta `build/docker/local/`. Você verá um arquivo chamado `.env-sample`. Faça uma cópia dele e renomeie a cópia para `.env`.
    ```bash
    # Navegue até a pasta
    cd build/docker/local/

    # Copie o arquivo de exemplo
    cp .env-sample .env
    ```
    *O arquivo `.env` já vem com as configurações padrão para o ambiente Docker local e não precisa de alterações para funcionar.*

3.  **Adicione um vídeo de teste:** Na raiz do projeto, navegue até a pasta `build/local_upload/` e coloque um arquivo de vídeo de exemplo. Renomeie-o para `trailerGTA6_4K.mp4` (ou o nome que estiver configurado no seu `docker-compose.infra.yml`).

### 3. Iniciando o Ambiente Completo

Com a configuração inicial pronta, iniciar todo o ambiente (infraestrutura + aplicação) é muito simples. Na **raiz do projeto**, execute um único comando:

```bash
make setup