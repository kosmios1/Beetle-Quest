import jwt
import sys
import uuid
import time
import pyotp
import string
import random
import logging
import binascii

from datetime import datetime
from http import HTTPStatus
from locust import HttpUser, task, FastHttpUser, between

base_path = "/api/v1"

user_ids = []
gacha_ids = []
auction_ids = []

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

        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

        response = self.client.post(f"{base_path}/auth/login", json={
            "username": self.username,
            "password": self.password,
        }, allow_redirects=False)

        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
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
        response = self.client.get(f"{base_path}/auth/logout", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

class UserMSRequests(AuthenticatedUser):
    wait_time = between(1, 2)

    @task
    def get_user(self):
        response = self.client.get(f"{base_path}/user/account/{self.user_id}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def update_user(self):
        response = self.client.patch(f"{base_path}/user/account/{self.user_id}", json={
            "username": "",
            "email": "",
            "new_password": self.password,
            "old_password": self.password,
        }, allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def delete_user(self):
        response = self.client.delete(f"{base_path}/user/account/delete", json={
            "password": self.password,
        }, allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

class GachaMSRequests(AuthenticatedUser):
    wait_time = between(1, 2)

    @task
    def get_gacha_list(self):
        response = self.client.get(f"{base_path}/gacha/list", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_gacha(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        response = self.client.get(f"{base_path}/gacha/{randgachaid}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_users_gacha_list(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        response = self.client.get(f"{base_path}/gacha/user/{randuserid}/list", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_users_gacha(self):
        if len(gacha_ids) == 0 or len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        randgachaid = random.choice(gacha_ids)
        response = self.client.get(f"{base_path}/gacha/user/{randgachaid}/{randuserid}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

class MarketMSRequests(AuthenticatedUser):
    wait_time = between(1, 2)
    @task
    def buy_bugscoin(self):
        response = self.client.post(f"{base_path}/market/bugscoin/buy", json={
            "amount": f"{100000}",
        }, allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def roll_gacha(self):
        response = self.client.get(f"{base_path}/market/gacha/roll", allow_redirects=False)

        if b'not enough money to roll gacha' in response.content:
            return
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def buy_gacha(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        response = self.client.get(f"{base_path}/market/gacha/{randgachaid}/buy", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def create_auction(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        response = self.client.post(f"{base_path}/market/auction", json={
            "gacha_id": randgachaid,
            "end_time": datetime.fromtimestamp(time.time() + 3600).strftime("%Y-%m-%dT%H:%M"),
        }, allow_redirects=False)

    @task
    def get_auction_list(self):
        response = self.client.get(f"{base_path}/market/auction/list", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_auction_details(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        response = self.client.get(f"{base_path}/market/auction/{randauctionid}", allow_redirects=False)

    @task
    def bid_to_auction(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        response = self.client.post(f"{base_path}/market/auction/{randauctionid}/bid", data={
            "bid_amount": f"{random.randint(0, 10000000)}",
        }, allow_redirects=False)

        resp_body = response.content
        good_err_msg = [ b'owner cannot bid', b'bid amount not enough', b'auction already ended', b'' ]
        for _, err in enumerate(good_err_msg):
            if err in resp_body:
                return

        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def delete_auction(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        response = self.client.delete(f"{base_path}/market/auction/{randauctionid}", data={
            "password": self.password,
        }, allow_redirects=False)

        resp_body = response.content
        good_err_msg = [ b'user not owner', b'invalid password', b'auction already ended', b'auction is too close to end', b'auction has bids']
        for _, err in enumerate(good_err_msg):
            if err in resp_body:
                return
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return


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
            "otp_code": self.otp.now()
        }, allow_redirects=False)

        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
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
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

class AdminMSRequests(AuthenticatedAdmin):
    wait_time = between(1, 2)
    """
    Login as an admin and then retrieve information to be used in the other tasks.
    """
    def on_start(self):
        global user_ids, gacha_ids, auction_ids
        super().on_start()

        response = self.client.get(f"{base_path}/admin/user/get_all", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

        response_data = response.json()
        user_list = response_data.get("UserList", [])
        user_ids = parse_uuids(user_list, "user_id")

        response = self.client.get(f"{base_path}/admin/gacha/get_all", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

        response_data = response.json()
        gacha_list = response_data.get("GachaList", [])
        gacha_ids = parse_uuids(gacha_list, "gacha_id")

        response = self.client.get(f"{base_path}/admin/market/auction/get_all", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

        response_data = response.json()
        auction_list = response_data.get("AuctionList", [])
        auction_ids = parse_uuids(auction_list, "auction_id")

    @task
    def get_user(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        response = self.client.get(f"{base_path}/admin/user/{randuserid}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def update_user(self):
        if len(user_ids) == 0:
            return
        randstr = generate_random_string()
        randuserid = random.choice(user_ids)
        response = self.client.patch(f"{base_path}/admin/user/{randuserid}", json={
            "username": randstr,
            "email": f"{randstr[:len(randstr)//2:]}@{randstr[len(randstr)//2::]}.it",
            "currency": f"{random.randint(0, 1000000)}",
        }, allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_user_transactions(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        response = self.client.get(f"{base_path}/admin/user/{randuserid}/transaction_history", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_user_auctions(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        response = self.client.get(f"{base_path}/admin/user/{randuserid}/auction/get_all", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
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
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_gacha_details(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        response = self.client.get(f"{base_path}/admin/gacha/{randgachaid}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def update_gacha(self):
        if len(gacha_ids) == 0:
            return
        randstr = generate_random_string()
        randgachaid = random.choice(gacha_ids)
        response = self.client.patch(f"{base_path}/admin/gacha/{randgachaid}", json={
            "name": randstr,
            "rarity": random.choice(["Common", "Uncommon", "Rare", "Epic", "Legendary"]),
            "price": f"{random.randint(0, 1000)}",
            "image_path": randstr[::-1]
        }, allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def delete_gacha(self):
        if len(gacha_ids) == 0:
            return
        randstr = generate_random_string()
        randgachaid = random.choice(gacha_ids)
        response = self.client.delete(f"{base_path}/admin/gacha/{randgachaid}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def get_transaction_history(self):
        response = self.client.get(f"{base_path}/admin/market/transaction_history", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return
    @task
    def get_auction_details(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        response = self.client.get(f"{base_path}/admin/market/auction/{randauctionid}", allow_redirects=False)
        if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
            response.raise_for_status()
            return

    @task
    def update_auction(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        randstr = generate_random_string()
        response = self.client.patch(f"{base_path}/admin/market/auction/{randauctionid}", allow_redirects=False)
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

def parse_uuids(data_list, uuid_field_name: str) -> list[str]:
    uuids = [str(
        uuid.UUID(bytes=bytes(data[uuid_field_name]))
    ) for data in data_list]
    return uuids
