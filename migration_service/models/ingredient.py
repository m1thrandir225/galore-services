import json


class Ingredient:
    def __init__(self, name, amount):
        self.name = name
        self.amount = amount

    def to_dict(self):
        return {"name": self.name, "amount": self.amount}