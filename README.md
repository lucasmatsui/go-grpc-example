# :fire: gRPC example
- Utilizado **GOlang**

### Oque foi implementado ?
- Um **server** gRPC
-  Criação de um serviço que :
	- Adiciona um usuário (**unary**)
	- Adiciona um usuário dando feedback de status de inserção em tempo real para o client (**server streaming**)
	- Adiciona vários usuários dando feedback de criação em tempo real para o server (**client streaming**)
	- Adiciona vários usuários mandando e recebendo os feedbacks de inserção (**client and server streaming**)
- Um **client** que faz as chamada gRPC
