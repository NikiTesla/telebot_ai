import telebot, openai
import os

telebot_key, open_api_key = os.getenv("TELEBOT_API"), os.getenv("OPENAI_API_KEY")
if telebot_key == "" or open_api_key == "":
    print("TELEBOT_API or OPENAI_API_KEY env is empty")
    exit(1)

bot = telebot.TeleBot(telebot_key)
openai.api_key = open_api_key

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