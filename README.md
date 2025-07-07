# Process Worker - Hackathon SOAT

Este projeto implementa o `hackthon-soat-process-worker`, worker responsável por processar vídeos de forma assíncrona para entrega do hackathon da pós graduação de Arquitetura de Software da FIAP.

O serviço opera em segundo plano, consumindo tarefas de uma fila de mensagens SQS. Para cada tarefa, ele baixa um vídeo de um serviço do S3, extrai todos os frames usando FFmpeg, compacta os frames em um único arquivo `.zip` e realiza o upload do arquivo em um S3 para download.

### ✨ Arquitetura

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

Siga estes passos para iniciar o ambiente de desenvolvimento completo na sua máquina (WIP)
