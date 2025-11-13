# Vita API - Sistema de Atendimento Domiciliar

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org)
[![Framework](https://img.shields.io/badge/framework-Gin-green.svg)](https://gin-gonic.com)
[![Database](https://img.shields.io/badge/database-MongoDB-green.svg)](https://www.mongodb.com)
[![Payments](https://img.shields.io/badge/payments-Stripe-6772E5.svg)](https://stripe.com)

## 📖 Visão Geral

A **Vita API** é o backend do sistema de atendimento domiciliar, uma plataforma de marketplace projetada para conectar Pacientes que necessitam de cuidados de enfermagem com Enfermeiros qualificados.

A API gerencia o ciclo de vida completo dos atendimentos, desde o cadastro e aprovação de profissionais, até a solicitação de visitas (agendadas e imediatas), processamento de pagamentos e comunicação em tempo real.

---

## 🚀 Principais Funcionalidades

* **Gerenciamento de Papéis:** Módulos distintos para Pacientes, Enfermeiros e Administradores.
* **Autenticação JWT:** Sistema seguro de autenticação e autorização baseado em tokens.
* **Sistema de Visitas:** Fluxo completo para solicitação, agendamento, confirmação e conclusão de visitas.
* **Pagamentos Integrados:** Integração com o **Stripe** para processamento de pagamentos (Payment Intents) e onboarding de enfermeiros (Stripe Connect).
* **Chat em Tempo Real:** Sistema de chat via **WebSocket** para comunicação direta entre pacientes e enfermeiros.
* **Aprovação de Cadastros:** Fluxo administrativo para aprovação de novos enfermeiros, incluindo upload de documentos.

---

## 💻 Stack Tecnológica

- **Linguagem:** Go (Golang)  
- **Framework Web:** Gin  
- **Banco de Dados:** MongoDB (usando mongo-driver oficial)  
- **Chat em Tempo Real:** Gorilla WebSocket  
- **Pagamentos:** Stripe  
- **Autenticação:** JSON Web Tokens (JWT)  
- **Documentação:** Swag (Swagger/OpenAPI)

---

## 📚 Documentação da API (Swagger)

A documentação completa e interativa de todos os endpoints está disponível via **Swagger UI**.

Após iniciar o servidor (veja instruções abaixo), acesse localmente:

> **[http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)**

### Para Desenvolvedores

Este projeto usa o `swag` para gerar a documentação OpenAPI a partir das anotações no código-fonte.

**Importante:** Se você modificar qualquer comentário de anotação (`// @Summary`, `// @Param`, etc.), você **deve** regenerar a documentação antes de comitar suas mudanças:

```bash
# Pare o servidor (se estiver rodando) e execute:
swag init
