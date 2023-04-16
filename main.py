import telebot
import os
import openai
import json

with open("config.json") as f:
    conf = json.load(f)
    print(conf)

bot = telebot.TeleBot(token=conf["TELEBOT_API"])
openai.api_key = conf["OPENAI_API_KEY"]

@bot.message_handler(func=lambda _: True)
def handle_message(message):
    response = openai.Completion.create(
        model="text-davinci-003",
        prompt=message.text,
        temperature=0.5,
        max_tokens=100,
        top_p=0.3,
        frequency_penalty=0.5,
        presence_penalty=0.0
    )

    bot.send_message(chat_id=message.from_user.id, text=response['choices'][0]['text'])

bot.polling()