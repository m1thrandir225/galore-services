import numpy as np
from transformers import BertTokenizer, BertModel
import torch


class TextEmbeddingModel:
    """TextEmbeddingModel class for generating embeddings for text data."""

    def __init__(self):
        self.tokenizer: BertTokenizer = BertTokenizer.from_pretrained(
            "bert-base-uncased"
        )  # pyright: ignore[reportUnknownMemberType]
        self.model: BertModel = BertModel.from_pretrained("bert-base-uncased")  # pyright: ignore[reportUnknownMemberType]

    def get_embeddings(self, text: str) -> np.ndarray:
        inputs = self.tokenizer(text, return_tensors="pt")
        outputs = self.model(**inputs)

        embeddings = outputs.last_hidden_state.mean(dim=1)
        return embeddings.detach().numpy()
