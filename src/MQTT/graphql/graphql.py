import os
import requests

queryLoginAccount = """
mutation LoginAccount($email: String!, $password: String!) {
    LoginAccount(account: { email: $email, password: $password }) {
        token
    }
}
"""

queryToggleDevice = """
  mutation ToggleDevice($id: Int!, $turnOn: Boolean!) {
    toggleDevice(device: { id: $id, turnOn: $turnOn })
  }
"""


class graphql:
  def __init__ (self):
    self.url = os.getenv("BACKEND_URL")
    self.email = os.getenv("EMAIL")
    self.password = os.getenv("PASSWORD")
    self.token = self.login
  def login(self):
    headers = {
      "Content-Type": "application/json"
    }
    variables = {
      "email": self.email,
      "password": self.password
    }
    response = requests.post(self.url, json={"query": queryLoginAccount, "variables": variables}, headers=headers)
    token = response.json()["data"]["LoginAccount"]["token"]
    print(token)
    return token
  def toggleDevice(self, id, turnOn):
    print(self.token)
    print(self.url)
    headers = {
      "Content-Type": "application/json",
      "Authorization": f"Bearer {self.token}"
    }
    variables = {
      "id": id,
      "turnOn": turnOn
    }
    print(variables)
    response = requests.post(self.url, json={"query": queryToggleDevice, "variables": variables}, headers=headers)
    if response.status_code == 401:
      self.token = self.login()
      headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {self.token}"
      }
      response = requests.post(self.url, json={"query": queryToggleDevice, "variables": variables}, headers=headers)
    return response.json()