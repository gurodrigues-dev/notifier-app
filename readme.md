# notifier-app

O notifier app, é uma aplicação resiliente e escalável feita para enviar alertas em diferentes tipos de canais, com facil integração e tolerante a falhas.

# Arquitetura

![Architecture](https://i.imgur.com/MeWpSLM.gif)

# Cronologia e pensamento

Inicialmente comecei construindo a arquitetura e tentando entender qual o melhor jeito de fazer as aplicações como um todo, não apenas no código. Mas principalmente na integração de novos canais e claro na escalabilidade.

Pensei em como seria receber essa requisição. Como não temos cenário de especificação se será uma api privada ou pública, se temos WAF, Api gateway, etc. Coloquei uma camada de segurança para validar se a requisição de fato, estaria ok, validando se o token existe no nosso banco de auth, em seguida. Adicionei uma validação no cache para toda vez que enviarmos uma mensagem via Kafka, ser adicionada no cache assim não haverá deduplicação da mensagem. Entretanto ainda era necessário enteder como o SLA poderia funcionar

Em seguida comecei a preparar os pacotes e foi me surgindo a dúvida. O que compensa mais fazer uma aplicação para notificar todos os tipos de canais ou criar uma aplicação/consumer para cada canal.

Levei em conta o tempo de configuração, deploy e a criação de cada aplicação. Então tive que começar a colocar na balança

Com um consumer eu criaria um God Object, quando um canal estivesse recebendo diversas mensagens. Todos escalariam em conjunto ao invés de um. A lógica de testes e a probilidade de quebra seria maior e a lógica provavelmente iria aumentar muito, não há um nível extremamente complexo, mas ainda assim aumentaria muito.

Entretanto, para adicionar um novo canal. Bastaria, escrever o código para conexão e adicionar no setup da aplicação, a leitura do kafka estaria pronta, o laço para enviar também. Então restaria apenas fazer o deploy e validar. Enquanto usar diversos, seria mais complexo, novas pipelines seriam necessárias, validar a pipeline. A depender de como a empresa pode fazer deploys. Pode ser demorado adicionar uma nova aplicação.

Pensando nisso, adicionei apenas uma aplicação. Para que o tempo de adicionar um novo canal fosse mais rápido, entretanto a escalabilidade iria perder um pouco da sua eficiência. Apesar de ambos os pontos serem essenciais nessa criação de aplicação. O que poderiamos fazer pra ajudar seria escalar verticalmente o consumer caso houvesse necessidade. Entendendo até mesmo o processamento de um consumer de notificações é pouca.

Com isso já imaginei que lidar com SLA, seria mais simples. Podemos criar um channel no Go, que aguarda erros de cada envio de mensagem. Caso todos falhem, podemos remover o cache do redis no final e avisar para um administrador que o envio de mensagem falhou.

Quando pensei nisso, instantaneamente imaginei que seria necessário um banco de dados para salvar as notificações, pelo menos momentaneamente. Por recebermos um json de notificação, esse json tem uma identificação de transação como uuid (eu assumi essa informação, não necessiaramente ela é verdadeira, mas levei isso em conta para minha decisão), poderíamos então pegar o json e com base no uuid salvar no banco.

Mas, talvez eu estaria adicionando itens demais ao redor desta aplicação eu poderia usar o próprio redis ao invés de usar SQL, ou mais um NoSQL. Isso adicionaria mais responsabilidade ao Redis, mas não mudaria seu papel principal. A única diferença é que poderiamos adicionar mais tempo de Cache, algo como 1 dia por exemplo.

Sei que o Redis não lida tão bem com tantas mensagens, então o tempo de 1 dia seria suficiente para que o mesmo não tivesse uma carga gigantesca de mensagens e poderíamos adicionar um campo de status da mensagem, algo como pending, error ou success. Assim quando eu visitar o cache ao receber a mensagem, caso o status for de sucesso, eu não reenvio a mensagem e se for de falha, envio novamente.

No status de falha eu poderia reenviar a mensagem novamente para o Kafka. Mas isso poderia causar um loop infinito, pois o erro que pode acontecer não necessariamente vai ser reparado rapidamente, logo pensei em adicionar mais um campo no json, como algum campo chamado de retries. Caso esse número seja superior a 3. Posso deixar ele lá no Redis parado durante um dia, em conjunto posso enviar para um banco de dados SQL, ou NoSQL e avisar algum administrador por alerta, com alertManager, PagerDuty, Sentry, algo similar e assim ele pode reenviar a mensagem ou entender o erro na aplicação. Obviamente, para reenviar a mensagem ele precisar ter um token no banco.

Com isso tive a ideia de

Criar uma rota de novos tokens e a tabela teria colunas similares a esta, porém para criar uma nova conta é necessário ter um usuário root que será configurado inicialmente na aplicação.
| id | token | admin_user

Para facilitar poderiamos ter uma rota para verificar o json enviado.
| id | uuid | json

Começando a construir mais a aplicação, me veio um questionamento. Além de ser fácil adicionar um novo ponto como um canal de um serviço diferente. Afinal, no slack, discord, email não vai existir apenas um canal serão vários.

No datadog, você separa os canais do slack por grupos, então vou enviar as notificações de maneira similar. Criarei uma nova tabela chamada de targets.
| id | platform | targetId | groupname |

Assim ao invés de chamar channels, vamos receber uma lista de targetIDs (que precisam estar registrados), assim enviaremos para todos os pontos. Em conjunto podemos receber os grupos.

# Rotas

## `api/v1/notify` (POST)
> Input
```json
{
  "uuid": "c3a4b7e0-1234-4d5f-abc0-56789abcd123",
  "message": "Você recebeu um pagamento de R$ 90,00!",
  "groups": ["vendas", "marketing"],
  "targetIds": ["slack-hook-123", "discord-hook-999"],
  "event": {
    "name": "payment_success",
    "timestamp": "1748190489",
    "cost_cents": "90000",
    "currency": "BRL",
    "requester": "requester-123@gmail.com",
    "receiver": "guester-123@gmail.com",
    "category": "pix"
  }
}
```

> Output
```json
{
  "message": "Sua mensagem <uuid> foi enviada com sucesso."
}
```
---

## `api/v1/token` (POST)
```json
{
  "admin_user": "James Sergio",
  "key": "environemt-root-key",
}
```

> Output
```json
{
  "message": "Token created"
}
```
---

## `api/v1/:uuid` (GET)
> Output
```json
{
  "uuid": "c3a4b7e0-1234-4d5f-abc0-56789abcd123",
  "channels": ["email", "slack"],
  "event": {
    "name": "payment_success",
    "timestamp": "1748190489",
    "cost_cents": "90000",
    "currency": "BRL",
    "requester": "requester-123@gmail.com",
    "receiver": "guester-123@gmail.com",
    "category": "pix"
  },
  "metadata": {
    "retries": 3,
    "status": "error", // error, success
    "error": "discord failed error, slack failed error, email failed error", // if sucess error must be null
  }
}
```





