import json


class Ingredient:
    def __init__(self, name, amount):
        self.name = name
        self.amount = amount

    def to_dict(self):
        return {"name": self.name, "amount": self.amount}

    def to_json(self):
        return json.dumps(self.to_dict())


class IngredientEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, Ingredient):
            return obj.__dict__  # or obj.to_dict() if you have a to_dict() method
        return super().default(obj)
