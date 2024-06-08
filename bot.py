import model
from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes
from typing import Final

start_command_text: Final = """Welcome to AI Bot, here you can talk with one of LLMs"""

class Bot:
    def __init__(self, telegram_key, bot_username):
        print("Starting bot...")
        if telegram_key == "":
            raise Exception("Telegram key is empty")
        self.bot_username = bot_username
        self.app = Application.builder().token(telegram_key).build()
        self.register_handlers()
        self.model = model.Model()
    
    def start(self):
        print("Polling...")
        self.app.run_polling()
    
    def register_handlers(self):
        self.app.add_handler(CommandHandler("start", self.start_command))
        self.app.add_handler(MessageHandler(filters.TEXT, self.handle_message))
        self.app.add_error_handler(self.error)

    async def start_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        await update.message.reply_text(start_command_text)

    async def handle_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        message_type: str = update.message.chat.type
        text: str = update.message.text
        print(f"User {update.message.chat.id} in {message_type}: '{text}'")
    
        if message_type == "group":
            print(f"group chat {text}")
            if self.bot_username in text:
                new_text: str = text.replace(self.bot_username, "").strip()
                response: str = self.handle_response(new_text)
            else:
                return
        else:
            response: str = self.handle_response(text)
        
        await update.message.reply_text(response)
    
    def handle_response(self, text: str) -> str:
        return self.model.discuss(text)
    
    async def error(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        print(f"Update {update} caused error: {context.error}")
