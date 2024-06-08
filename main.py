import bot
import os

if __name__ == "__main__":
    telegram_key = os.getenv("TELEGRAM_API_KEY")
    bot_username = os.getenv("BOT_USERNAME")
    if telegram_key == "":
        print("TELEGRAM_API_KEY env is empty")
        exit(1)
    
    telebot = bot.Bot(telegram_key, bot_username)
    telebot.start()
