### Exposing public URLs for your local webserver.
https://ngrok.com/

**Start ngrok**

`ngrok http <port>`

### How to configure Telegram Webhook
`curl -F "url=https://<ngrok_id>.ngrok.io/update"  https://api.telegram.org/bot<your_api_token>/setWebhook`
