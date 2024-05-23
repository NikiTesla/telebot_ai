import telebot
import openai
import os

telebot_key, open_api_key = os.getenv("TELEBOT_API"), os.getenv("OPENAI_API_KEY")
if telebot_key == "" or open_api_key == "":
    print("TELEBOT_API or OPENAI_API_KEY env is empty")
    exit(1)

client = openai.OpenAI(
    api_key=open_api_key,
)
bot = telebot.TeleBot(telebot_key)

@bot.message_handler(func=lambda _: True)
def handle_message(message):
    response = client.chat.completions.create(
        model = "gpt-3.5-turbo-16k",
        messages = [message.text],
    )

    bot.send_message(chat_id=message.from_user.id, text=response['choices'][0]['text'])

bot.polling()