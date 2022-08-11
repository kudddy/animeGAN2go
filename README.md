## DialogFlow

Easy platform who care about your support in telegram.


### Contracts
Url format for tlg webhook.
```
http://lolka.ru/d5f53-n7ks8-new020-as322f/bot
```
where `d5f53-n7ks8-new020-as322f` - project id[auth parameter], `bot` - customer type[operator or bot].


Example url if you want to create project.
```
http://lolka.ru/b1630dbc-51a4-4462-81c8-5233d2a92081/update
```
where `b1630dbc-51a4-4462-81c8-5233d2a92081` - your auth id, `update` - it is route
If you want to create project, you must create tokes for bot-user and bot-operator and create bot in SmartMarket and take weebhook
Example payload:
```
{
    "bot": "***",
    "operator": "***",
    "sm-webhook": "https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/cGnGPZWb:45c9c4e54edfcf2cfe505f84e3f338185a334e42"
}
```
generate token

