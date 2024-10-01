import json
import os
from os.path import isfile
import requests


class Parser:
    def __init__(self, url: str):
        self.url = url

    # TODO: rewrite this, this will be a cron job I don't know if it is smart to
    # be saved in a file
    def get_cocktails(self):
        response = requests.get(self.url)

        json_data = response.json()["drinks"]  # the list of drinks

        # save the list to a file
        with open("cocktails_by_category.json", "w") as output:
            output.write(json_data)

    def openFile(self):
        # check if file exists
        if os.path.isfile("./cocktails_by_category.json") is False:
            self.get_cocktails()

        file = open("cocktails_by_category.json", "r")

        data = json.load(file)

        # TODO: parse the content

        file.close()
