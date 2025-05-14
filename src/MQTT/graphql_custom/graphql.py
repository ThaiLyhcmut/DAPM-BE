import os
import requests
import asyncio
from gql import Client, gql
from gql.transport.websockets import WebsocketsTransport

# Giữ GraphQL Queries dưới dạng chuỗi cho HTTP requests
queryGetHome_str = """
query GetHome {
  getHome {
    area {
      equipment {
        id
        title
        timeStart
        timeEnd
        turnOn
        cycle
      }
    }
  }
}
"""

queryLoginAccount_str = """
mutation LoginAccount($email: String!, $password: String!) {
  LoginAccount(account: { email: $email, password: $password }) {
    token
  }
}
"""

queryToggleDevice_str = """
mutation ToggleDevice($id: Int!, $turnOn: Boolean!) {
  toggleDevice(device: { id: $id, turnOn: $turnOn })
}
"""

# Subscription query sẽ được chuyển đổi thành DocumentNode trong phương thức subscribe
subscriptionQuery_str = """
subscription DeviceStatusUpdated {
  deviceStatusUpdated {
    id
    turnOn
  }
}
"""

class GraphQL:
    def __init__(self):
        self.url = os.getenv("BACKEND_URL")  # e.g. https://your-api/graphql
        self.ws_url = os.getenv("BACKEND_WS_URL")  # e.g. wss://your-api/graphql
        # Nếu BACKEND_WS_URL không được đặt, tự động chuyển đổi từ BACKEND_URL
        if not self.ws_url and self.url:
            self.ws_url = self.url.replace("http://", "ws://").replace("https://", "wss://")
        self.email = os.getenv("EMAIL")
        self.password = os.getenv("PASSWORD")
        self.token = self.login()
        self.name = os.getenv("MQTT_NAME")

    def login(self):
        headers = {
            "Content-Type": "application/json"
        }
        variables = {
            "email": self.email,
            "password": self.password
        }
        response = requests.post(self.url, json={"query": queryLoginAccount_str, "variables": variables}, headers=headers)
        token = response.json()["data"]["LoginAccount"]["token"]
        print("✅ Logged in:", token)
        return token

    def get_home(self):
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.token}"
        }
        response = requests.post(self.url, json={"query": queryGetHome_str}, headers=headers)
        if response.status_code == 401:
            self.token = self.login()
            headers["Authorization"] = f"Bearer {self.token}"
            response = requests.post(self.url, json={"query": queryGetHome_str}, headers=headers)
        data = response.json()
        # Xử lý dữ liệu để tạo danh sách chuỗi title_id
        result = []
        for home in data["data"]["getHome"]:
            for area in home["area"]:
                if (area["equipment"] != None):
                    for equipment in area["equipment"]:
                        result.append(f"{self.name}{equipment['title']}_{equipment['id']}")
        return result

    def toggleDevice(self, id, turnOn):
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.token}"
        }
        variables = {
            "id": id,
            "turnOn": turnOn
        }
        response = requests.post(self.url, json={"query": queryToggleDevice_str, "variables": variables}, headers=headers)
        if response.status_code == 401:
            self.token = self.login()
            headers["Authorization"] = f"Bearer {self.token}"
            response = requests.post(self.url, json={"query": queryToggleDevice_str, "variables": variables}, headers=headers)
        return response.json()

    async def subscribe(self, client):
        # Chuyển đổi chuỗi query thành DocumentNode chỉ khi cần thiết
        subscription_document = gql(subscriptionQuery_str)
    
        # Sử dụng tham số headers thay vì init_payload
        headers = {"Authorization": f"Bearer {self.token}"}
        
        transport = WebsocketsTransport(
            url=self.ws_url,
            headers=headers  # Truyền token qua headers
    )
        
        async with Client(transport=transport, fetch_schema_from_transport=True) as session:
            async for result in session.subscribe(subscription_document):
                print(result["deviceStatusUpdated"])
                yield result["deviceStatusUpdated"]