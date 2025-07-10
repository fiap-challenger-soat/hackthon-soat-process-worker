# hackthon-soat-process-worker

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=fiap-challenger-soat_hackthon-soat-process-worker&metric=alert_status&token=5000fa51b16d2712a14dda6e6ff3474180a2a6ed)](https://sonarcloud.io/summary/new_code?id=fiap-challenger-soat_hackthon-soat-process-worker)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=fiap-challenger-soat_hackthon-soat-process-worker&metric=coverage&token=5000fa51b16d2712a14dda6e6ff3474180a2a6ed)](https://sonarcloud.io/summary/new_code?id=fiap-challenger-soat_hackthon-soat-process-worker)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=fiap-challenger-soat_hackthon-soat-process-worker&metric=security_rating&token=5000fa51b16d2712a14dda6e6ff3474180a2a6ed)](https://sonarcloud.io/summary/new_code?id=fiap-challenger-soat_hackthon-soat-process-worker)

---

## ✨ Visão Geral

O **hackthon-soat-process-worker** é um serviço backend assíncrono para processamento automático de vídeos. Ele consome tarefas de uma fila SQS, baixa vídeos de um bucket S3, extrai frames usando FFmpeg, compacta os frames em um arquivo `.zip` e disponibiliza o resultado para download. Todo o fluxo é automatizado e simulado localmente via LocalStack e Docker.

---

### 📐 System Design

#### Caso de Sucesso

![Caso de Sucesso]([https://drive.google.com/drive/u/0/folders/1-NE70aPGQ095kNQRzVqoVaI8tMJSww23](https://drive.google.com/file/d/1xbmCnDGWXq3DSND3PnkEtpj2DmeMOTIx/view?usp=sharing))

#### Caso de Erro

[Visualizar Caso de Erro (Google Drive)]([https://drive.google.com/drive/u/0/folders/1-NE70aPGQ095kNQRzVqoVaI8tMJSww23](https://drive.google.com/file/d/1iJxWQd4J48whN1iyDwluCrOIx4qnE4RZ/view?usp=sharing))

---

### 🛠️ Pré-requisitos

Certifique-se de ter instalado:

- **Git** – Para clonar o repositório.
- **Go** – Versão 1.22 ou superior.
- **Docker & Docker Desktop** – Para executar a infraestrutura em contêineres.
- **Make** – Para facilitar comandos via `Makefile`.
  - **Windows:** Use o terminal **Git Bash** (já inclui o `make`) ou instale via Chocolatey: `choco install make`.
- **AWS CLI** – Para interações manuais com o LocalStack.
- **FFmpeg** – Utilizado no processamento de vídeos localmente.

---

### 🚀 Como Rodar

#### 1. Adicione arquivos de vídeo e imagem

Inclua um vídeo e um arquivo de imagem na pasta `build/local_upload`. Isso simula duas mensagens SQS: uma para processamento normal e outra para erro.

#### 2. Configure os nomes dos arquivos

- No arquivo `build/docker/local/init-aws.sh`, altere o array de arquivos (linha 13) substituindo `"your-video-file"` pelo nome dos arquivos que você adicionou.
- No `build/docker/local/docker-compose.yml`, ajuste os volumes do container LocalStack, substituindo `"your-video-file"` pelos nomes reais dos arquivos.

#### 3. Suba os containers

```sh
make up
```

Toda a infraestrutura será inicializada e o fluxo de processamento começará automaticamente.

#### 4. Execute o projeto

Com os containers rodando, acompanhe os logs para visualizar todas as etapas do processamento.

Para reiniciar o fluxo:

```sh
make down
make up
```

---

### 🗄️ Migrations e Seeding

O projeto inclui migrations para o banco Postgres, simulando a conexão e o seeding de dados necessários para o funcionamento do fluxo.

---

### 🧰 Comandos Úteis para AWS LocalStack

- **Verificar buckets existentes**
  ```sh
  aws --endpoint-url=http://localhost:4566 s3 ls
  ```

- **Listar conteúdo do bucket de saída**
  ```sh
  aws --endpoint-url=http://localhost:4566 s3 ls s3://bucket-videos/output/
  ```

- **Fazer download do arquivo processado**
  ```sh
  aws --endpoint-url=http://localhost:4566 s3 cp s3://bucket-videos/output/NOME_DO_ARQUIVO.zip .
  ```

- **Verificar quantidade de mensagens na fila de erro**
  ```sh
  aws --endpoint-url=http://localhost:4566 sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/error-queue --attribute-names ApproximateNumberOfMessages
  ```

- **Verificar mensagens da fila de erro**
  ```sh
  aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/error-queue
  ```

---

### 💡 Observações

- O script `build/docker/local/init-aws.sh` automatiza toda a configuração dos recursos AWS simulados e dispara as mensagens SQS para o worker.
- O processamento é feito automaticamente ao subir os containers.
- Os logs detalham cada etapa do processamento, incluindo erros e sucesso.
