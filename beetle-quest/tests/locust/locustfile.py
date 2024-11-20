import jwt
import sys
import pyotp
import string
import random
import logging
import binascii

from http import HTTPStatus
from locust import HttpUser, task, FastHttpUser

base_path = "/api/v1"

class AuthenticatedUser(FastHttpUser):
    abstract = True
    insecure = True

    user_id = None
    username = None
    password = None
    email = None

    access_token = None

    """
    This method is called before the virtual user execute any task.
    It is used to perform the login and registration of a random user and to store
    its UUID.
    """
    def on_start(self):
        random_str = generate_random_string()
        self.username = random_str
        self.password = random_str[::-1]
        self.email = f"{random_str[:len(random_str)//2]}@{random_str[len(random_str)//2:]}.it"

        response = self.client.post(f"{base_path}/auth/register", json={
            "username": self.username,
            "password": self.password,
            "email": self.email,
        }, allow_redirects=False)

        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

        response = self.client.post(f"{base_path}/auth/login", json={
            "username": self.username,
            "password": self.password,
        }, allow_redirects=False)

        if response.status_code != HTTPStatus.FOUND:
            response.raise_for_status()
            return

        str_tok = None
        for cookie in self.client.cookiejar:
            if cookie.name == "access_token":
                str_tok = cookie.value
                break

        if str_tok == None:
            print("Failed to retrive the token from the login response!")
            sys.exit();

        self.access_token = parse_jwt(str_tok)
        if self.access_token == None:
            print(f"Failed to parse the token {str_tok}")
            sys.exit()
        self.user_id = self.access_token["sub"]

    def on_stop(self):
        response = self.client.delete(f"{base_path}/user/account/delete", json={
            "password": self.password,
        })
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

class UserMSRequests(AuthenticatedUser):
    @task
    def get_user(self):
        response = self.client.get(f"{base_path}/user/account/{self.user_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def update_user(self):
        response = self.client.put(f"{base_path}/user/account/{self.user_id}", json={
            "email": self.email,
            "password": self.password,
        }, allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    # NOTE: delete user is called when inside the on_stop function of the AuthenticatedUser class

class GachaMSRequests(AuthenticatedUser):
    @task
    def get_gacha_list(self):
        response = self.client.get(f"{base_path}/gacha", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_gacha(self):
        response = self.client.get(f"{base_path}/gacha/{self.gacha_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_users_gacha_list(self):
        response = self.client.get(f"{base_path}/gacha/{self.user_id}/list", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_users_gacha(self):
        response = self.client.get(f"{base_path}/{self.gacha_id}/{self.user_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

# TODO: class MarketMSRequests(AuthenticatedUser):

# ==============================================================================
# Admin
# ==============================================================================

class AuthenticatedAdmin(FastHttpUser):
    abstract = True
    insecure = True

    # NOTE: the admin parameters are hardcoded for simplicity reasons.
    # They are present in the database and are used to login the admin every
    # time the server goes up. A real implementation would require a more
    # sophisticated approach.
    admin_id = "09087f45-5209-4efa-85bd-761562a6df53"
    password = "admin"
    email = "admin@admin.com"
    otp = pyotp.TOTP("g2ytwh764px5wzorxcbk2c2f2jhv74kd")

    access_token = None

    def on_start(self):
        response = self.client.post(f"{base_path}/auth/admin/login", json={
            "admin_id": self.admin_id,
            "password": self.password,
            "otp_code": "123456",
        }, allow_redirects=False)

        if response.status_code != HTTPStatus.FOUND:
            response.raise_for_status()
            return

        str_tok = None
        for cookie in self.client.cookiejar:
            if cookie.name == "access_token":
                str_tok = cookie.value
                break

        if str_tok == None:
            print("Failed to retrive the token from the login response!")
            sys.exit();

        self.access_token = parse_jwt(str_tok)
        if self.access_token == None:
            print(f"Failed to parse the token {str_tok}")
            sys.exit()

    def on_stop(self):
        response = self.client.get(f"{base_path}/auth/logout")
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

class AdminMSRequests(AuthenticatedAdmin):
    @task
    def get_users(self):
        response = self.client.get(f"{base_path}/admin/user/get_all", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_user(self):
        response = self.client.get(f"{base_path}/admin/user/{self.user_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def update_user(self):
        randstr = generate_random_string()
        response = self.client.patch(f"{base_path}/admin/user/{self.user_id}", json={
            "username": randstr,
            "email": f"{randstr[:len(randstr)//2:]}@{randstr[len(randstr)//2::]}.it",
            "currency": f"{random.randint(0, 1000)}",
        }, allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_user_transactions(self):
        response = self.client.get(f"{base_path}/admin/user/{self.user_id}/transaction_history", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_user_auctions(self):
        response = self.client.get(f"{base_path}/admin/user/{self.user_id}/get_all", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def add_gacha(self):
        randstr = generate_random_string()
        response = self.client.post(f"{base_path}/admin/gacha/add", json={
            "name": randstr,
            "rarity": random.choice(["Common", "Uncommon", "Rare", "Epic", "Legendary"]),
            "price": f"{random.randint(0, 1000)}",
            "image_path": randstr[::-1]
        }, allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_all_gacha(self):
        response = self.client.get(f"{base_path}/admin/gacha/get_all", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_gacha_details(self):
        response = self.client.get(f"{base_path}/admin/gacha/{self.gacha_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def update_gacha(self):
        randstr = generate_random_string()
        response = self.client.patch(f"{base_path}/admin/gacha/{self.gacha_id}", json={
            "name": randstr,
            "rarity": random.choice(["Common", "Uncommon", "Rare", "Epic", "Legendary"]),
            "price": f"{random.randint(0, 1000)}",
            "image_path": randstr[::-1]
        }, allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def delete_gacha(self):
        randstr = generate_random_string()
        response = self.client.delete(f"{base_path}/admin/gacha/{self.gacha_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_transaction_history(self):
        response = self.client.get(f"{base_path}/admin/market/transaction_history", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return
    @task
    def get_all_auctions(self):
        response = self.client.get(f"{base_path}/admin/market/auction/get_all", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def get_auction_details(self):
        response = self.client.get(f"{base_path}/admin/market/auction/{self.auction_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.OK:
            response.raise_for_status()
            return

    @task
    def update_auction(self):
        randstr = generate_random_string()
        response = self.client.patch(f"{base_path}/admin/market/auction/{self.auction_id}", allow_redirects=False)
        if response.status_code != HTTPStatus.NOT_IMPLEMENTED: # TODO: Change to OK
            response.raise_for_status()
            return

# ==============================================================================
# Utils
# ==============================================================================

def generate_random_string(length=10):
    characters = string.ascii_letters + string.digits
    random_string = ''.join(random.choices(characters, k=length))
    return random_string

def parse_jwt(token: str, algorithms=["HS256"]) -> dict | None:
    try:
        decoded_token = jwt.decode(token, bytes(), algorithms=algorithms, options={"verify_signature": False})
        return decoded_token
    except Exception:
        return None
