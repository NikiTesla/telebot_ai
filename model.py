from transformers import pipeline

class Model:
    def __init__(self, model="facebook/blenderbot-400M-distill"):
        self.chatbot = pipeline(task="conversational", model=model)
    
    def discuss(self, prompt: str) -> str:
        conversation = self.chatbot(prompt)
        return conversation[-1]["generated_text"]