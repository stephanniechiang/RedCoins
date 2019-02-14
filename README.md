# RedCoins

Para executar a API antes deve-se antes instalar os seguintes packages:
- gorilla/mux
- jinzhu/gorm
- dgrijalva/jwt-go
- joho/godotenv

Utilizando o comando: 'go get github.com/{nome do package}' no terminal e criar um banco de dados vazio Postgres chamado 'redcoins'.
Após isso para executar a API deve-se rodar o seguinte comando 'go build main.go' e em seguida 'go ./main'.
As seguintes rotas fazem parte da API:
#### Cria um novo usuário
##### POST: /api/user/new
		Estrutura (exemplo):
			{"name": "João",
			"last_name": "Peixoto",
    		"email": "joao@peixoto.com",
    		"password": "12345678",
			"birthday": "11/10/1994",
			"balance": 0}
#### Autentica o usuário
##### POST: /api/user/login
		Estrutura (exemplo):
			{ "email": "joao@peixoto.com",
			"password": "12345678"}
#### Cria nova transação de venda
##### POST: /api/transactions/newsell
		Token do usuário: obtido no login
		Estrutura (exemplo):
			{"bitcoins": 321,
			"user_id_2": 1}
#### Cria nova transação de compra
##### POST: /api/transactions/newbuy
		Token do usuário: obtido no login
		Estrutura (exemplo):
			{"bitcoins": 120,
			"user_id_2": 1}
#### Obtem transações do usuário autenticado
##### GET: /api/user/me/transactions
		Token do usuário: obtido no login
#### Obtem transações do usuário escolhido
##### GET: /api/user/{userId}/transactions
		Token do usuário: obtido no login
