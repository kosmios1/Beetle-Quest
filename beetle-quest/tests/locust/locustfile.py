import jwt
import sys
import uuid
import time
import pyotp
import secrets
import random

import hashlib
import base64
import traceback
import threading
import logging

from datetime import datetime
from http import HTTPStatus
from locust import task, between, FastHttpUser, LoadTestShape
from bs4 import BeautifulSoup


base_path = "/api/v1"

user_ids = []
gacha_ids = []
auction_ids = []

lock = threading.Lock()
total_counter = 0
rarities_counter = {}
rarities_percentage = {}

class AuthenticatedUser(FastHttpUser):
    abstract = True
    insecure = True

    user_id = None
    username = None
    password = None
    email = None

    access_token = None
    identity_token = None

    host = "https://localhost"

    def __init__(self, environment) -> None:
        self.host = "https://localhost"
        super().__init__(environment)

    def on_start(self):
        self.make_authentication_request()
        self.make_oauth_authorization_request()

    def make_authentication_request(self):
        random_str = generate_random_string()
        self.username = random_str
        self.password = random_str[::-1]
        self.email = f"{random_str[:len(random_str)//2]}@{random_str[len(random_str)//2:]}.it"

        with self.client.post(f"{base_path}/auth/register", json={
            "username": self.username,
            "password": self.password,
            "email": self.email,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

        with self.client.post(f"{base_path}/auth/login", json={
            "username": self.username,
            "password": self.password,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    def make_oauth_authorization_request(self):
        code = ""
        state = generate_random_string()
        codeVerifier, codeChallenge = generate_code_verifier_and_chall()
        with self.client.get("/oauth/authorize", params={
            "response_type": "code",
            "client_id": "beetle-quest",
            "redirect_uri": "/fake/callback",
            "scope": "gacha, user, market",
            "state": state,
            "code_challenge": codeChallenge,
            "code_challenge_method": "S256",
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return

            if response.headers is not None:
                location = response.headers["Location"]
                if location and len(location.split("code=")) > 1:
                    code = location.split("code=")[1].split("&")[0]

            response.success()

        with self.client.post("/oauth/token", data={
            "grant_type": "authorization_code",
            "code": code,
            "redirect_uri": "/fake/callback",
            "client_id": "beetle-quest",
            "code_verifier": codeVerifier,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return

            if response.text is not None:
                if "access_token" in response.json() and "id_token" in response.json():
                    self.access_token = response.json()["access_token"]
                    self.client.auth_header = f"Bearer {self.access_token}"

                    self.identity_token = response.json()["id_token"]
                    if self.identity_token == None:
                        print(f"Failed to parse the token {self.identity_token}")
                        return

                    id_token = parse_jwt(self.identity_token, algorithms="HS256")
                    if id_token == None:
                        print(f"Failed to parse the token {id_token}")
                        return
                    self.user_id = id_token["sub"]
            response.success()

    def on_stop(self):
        with self.client.get(f"{base_path}/auth/logout", catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

class UserMSRequests(AuthenticatedUser):
    wait_time = between(1, 2)

    @task
    def get_userinfo(self):
        with self.client.get("/userinfo", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_user(self):
        with self.client.get(f"{base_path}/user/account/{self.user_id}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def update_user(self):
        with self.client.patch(f"{base_path}/user/account/{self.user_id}", json={
            "username": "",
            "email": "",
            "new_password": self.password,
            "old_password": self.password,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def delete_user(self):
        if random.random() >= 0.1:
            return

        with self.client.post(f"{base_path}/user/account/{self.user_id}", json={
            "password": self.password,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

class GachaMSRequests(AuthenticatedUser):
    wait_time = between(1, 2)

    @task
    def get_gacha_list(self):
        with self.client.get(f"{base_path}/gacha/list", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_gacha(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        with self.client.get(f"{base_path}/gacha/{randgachaid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_users_gacha_list(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        with self.client.get(f"{base_path}/gacha/user/{randuserid}/list", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_users_gacha(self):
        if len(gacha_ids) == 0 or len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        randgachaid = random.choice(gacha_ids)
        with self.client.get(f"{base_path}/gacha/{randgachaid}/user/{randuserid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

class MarketMSRequests(AuthenticatedUser):
    own_auctions = []
    own_gacha = []

    wait_time = between(1, 2)
    @task
    def buy_bugscoin(self):
        with self.client.post(f"{base_path}/market/bugscoin/buy", json={
            "amount": f"{100000}",
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def roll_gacha(self):
        global total_counter, rarities_counter, rarities_percentage

        with self.client.get(f"{base_path}/market/gacha/roll", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            if response.text != None:
                soup = BeautifulSoup(response.text, 'html.parser')
                input_tag = soup.find('input', {'name': 'hidden_data'})
                if input_tag:
                    hidden_data_value = input_tag["value"]
                    self.own_gacha.append(hidden_data_value)

                input_tag = soup.find('input', {'name': 'hidden_data2'})
                if input_tag:
                    hidden_data_value = input_tag["value"]
                    if hidden_data_value != '':
                        with lock:
                            if hidden_data_value not in rarities_counter:
                                rarities_counter[hidden_data_value] = 0
                            total_counter = total_counter + 1
                            rarities_counter[hidden_data_value] = rarities_counter[hidden_data_value] + 1
                            for key in rarities_counter:
                                rarities_percentage[key] = (rarities_counter[key] / total_counter) * 100
                        logging.info(rarities_percentage)
            response.success()

    @task
    def buy_gacha(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        with self.client.get(f"{base_path}/market/gacha/{randgachaid}/buy", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            if response.text != None:
                soup = BeautifulSoup(response.text, 'html.parser')
                input_tag = soup.find('input', {'name': 'hidden_data'})
                if input_tag:
                    hidden_data_value = input_tag["value"]
                    self.own_gacha.append(hidden_data_value)
            response.success()

    @task
    def create_auction(self):
        if len(self.own_gacha) == 0:
            return
        randgachaid = random.choice(self.own_gacha)
        with self.client.post(f"{base_path}/market/auction/", json={
            "gacha_id": randgachaid,
            "end_time": datetime.fromtimestamp(time.time() + (5 * 60)).strftime("%Y-%m-%dT%H:%M"),
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            if response.text != None:
                soup = BeautifulSoup(response.text, 'html.parser')
                input_tag = soup.find('input', {'name': 'hidden_data'})
                if input_tag:
                    hidden_data_value = input_tag["value"]
                    self.own_auctions.append(hidden_data_value)
            response.success()

    @task
    def get_auction_list(self):
        with self.client.get(f"{base_path}/market/auction/list", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_auction_details(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        with self.client.get(f"{base_path}/market/auction/{randauctionid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def bid_to_auction(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        with self.client.post(f"{base_path}/market/auction/{randauctionid}/bid", json={
            "bid_amount": f"{random.randint(0, 100000)}",
            }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def delete_auction(self):
        if random.random() >= 0.1:
            return

        if len(self.own_auctions) == 0:
            return
        randauctionid = random.choice(self.own_auctions)
        with self.client.post(f"{base_path}/market/auction/{randauctionid}", json={
            "password": self.password,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

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
    identity_token = None

    def __init__(self, environment) -> None:
        self.host = "https://localhost:6443"
        super().__init__(environment)

    def on_start(self):
        self.make_authentication_request()
        self.make_oauth_authorization_request()

    def make_authentication_request(self):
        with self.client.post(f"{base_path}/auth/admin/login", json={
            "admin_id": self.admin_id,
            "password": self.password,
            "otp_code": self.otp.now()
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    def make_oauth_authorization_request(self):
        code = ""
        state = generate_random_string()
        codeVerifier, codeChallenge = generate_code_verifier_and_chall()
        with self.client.get("/oauth/authorize", params={
            "response_type": "code",
            "client_id": "beetle-quest",
            "redirect_uri": "/fake/callback",
            "scope": "admin",
            "state": state,
            "code_challenge": codeChallenge,
            "code_challenge_method": "S256",
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return

            if response.headers is not None:
                location = response.headers["Location"]
                if location and len(location.split("code=")) > 1:
                    code = location.split("code=")[1].split("&")[0]

            response.success()

        with self.client.post("/oauth/token", data={
            "grant_type": "authorization_code",
            "code": code,
            "redirect_uri": "/fake/callback",
            "client_id": "beetle-quest",
            "code_verifier": codeVerifier,
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return

            if response.text is not None:
                if "access_token" in response.json() and "id_token" in response.json():
                    self.access_token = response.json()["access_token"]
                    self.client.auth_header = f"Bearer {self.access_token}"

                    self.identity_token = response.json()["id_token"]
                    if self.identity_token == None:
                        print(f"Failed to parse the token { self.identity_token}")
                        return

                    id_token = parse_jwt(self.identity_token, algorithms="HS256")
                    if id_token == None:
                        print(f"Failed to parse the token {id_token}")
                        return
                    self.user_id = id_token["sub"]
            response.success()

    def on_stop(self):
        with self.client.get(f"{base_path}/auth/logout", catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

class AdminMSRequests(AuthenticatedAdmin):
    wait_time = between(1, 2)
    """
    Login as an admin and then retrieve information to be used in the other tasks.
    """
    def on_start(self):
        global user_ids, gacha_ids, auction_ids
        super().on_start()

        with self.client.get(f"{base_path}/admin/user/get_all", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response_data = response.json()
            user_list = response_data.get("UserList", [])
            user_ids = parse_uuids(user_list, "user_id")

        with self.client.get(f"{base_path}/admin/gacha/get_all", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response_data = response.json()
            gacha_list = response_data.get("GachaList", [])
            gacha_ids = parse_uuids(gacha_list, "gacha_id")

        with self.client.get(f"{base_path}/admin/market/auction/get_all", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response_data = response.json()
            auction_list = response_data.get("AuctionList", [])
            auction_ids = parse_uuids(auction_list, "auction_id")

    @task
    def get_user(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        with self.client.get(f"{base_path}/admin/user/{randuserid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def update_user(self):
        if len(user_ids) == 0:
            return
        randstr = generate_random_string()
        randuserid = random.choice(user_ids)
        with self.client.patch(f"{base_path}/admin/user/{randuserid}", json={
            "username": randstr,
            "email": f"{randstr[:len(randstr)//2:]}@{randstr[len(randstr)//2::]}.it",
            "currency": f"{random.randint(0, 1000000)}",
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_user_transactions(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        with self.client.get(f"{base_path}/admin/user/{randuserid}/transaction_history", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_user_auctions(self):
        if len(user_ids) == 0:
            return
        randuserid = random.choice(user_ids)
        with self.client.get(f"{base_path}/admin/user/{randuserid}/auction/get_all", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def add_gacha(self):
        randstr = generate_random_string()
        with self.client.post(f"{base_path}/admin/gacha/add", json={
            "name": randstr,
            "rarity": random.choice(["Common", "Uncommon", "Rare", "Epic", "Legendary"]),
            "price": f"{random.randint(0, 1000)}",
            "image_path": randstr[::-1]
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_gacha_details(self):
        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        with self.client.get(f"{base_path}/admin/gacha/{randgachaid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def update_gacha(self):
        if len(gacha_ids) == 0:
            return
        randstr = generate_random_string()
        randgachaid = random.choice(gacha_ids)
        with self.client.patch(f"{base_path}/admin/gacha/{randgachaid}", json={
            "name": randstr,
            "rarity": random.choice(["Common", "Uncommon", "Rare", "Epic", "Legendary"]),
            "price": f"{random.randint(0, 1000)}",
            "image_path": randstr[::-1]
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def delete_gacha(self):
        if random.random() >= 0.1:
            return

        if len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        with self.client.delete(f"{base_path}/admin/gacha/{randgachaid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_transaction_history(self):
        with self.client.get(f"{base_path}/admin/market/transaction_history", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def get_auction_details(self):
        if len(auction_ids) == 0:
            return
        randauctionid = random.choice(auction_ids)
        with self.client.get(f"{base_path}/admin/market/auction/{randauctionid}", allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

    @task
    def update_auction(self):
        if len(auction_ids) == 0 or len(gacha_ids) == 0:
            return
        randgachaid = random.choice(gacha_ids)
        randauctionid = random.choice(auction_ids)
        with self.client.patch(f"{base_path}/admin/market/auction/{randauctionid}", json={
            "gacha_id": randgachaid
        }, allow_redirects=False, catch_response=True) as response:
            if response.status_code == HTTPStatus.INTERNAL_SERVER_ERROR:
                response.failure(f"Request failed with status code: {response.status_code}")
                return
            response.success()

# ==============================================================================
# Test stages definition
# ==============================================================================

class StagesShapeWithCustomUsers(LoadTestShape):
    stages = [
        {"duration": 15, "users": 10, "spawn_rate": 5, "user_classes": [UserMSRequests, GachaMSRequests]},
        {"duration": 45, "users": 30, "spawn_rate": 10, "user_classes": [MarketMSRequests]},
        {"duration": 75, "users": 50, "spawn_rate": 15, "user_classes": [AdminMSRequests, UserMSRequests, GachaMSRequests]},
        {"duration": 150, "users": 70, "spawn_rate": 25, "user_classes": [MarketMSRequests]},
        {"duration": 270, "users": 100, "spawn_rate": 30, "user_classes": [AdminMSRequests, UserMSRequests, GachaMSRequests, MarketMSRequests]},
        {"duration": 345, "users": 70, "spawn_rate": 25, "user_classes": [MarketMSRequests]},
        {"duration": 390, "users": 50, "spawn_rate": 15, "user_classes": [AdminMSRequests, UserMSRequests, GachaMSRequests]},
        {"duration": 420, "users": 30, "spawn_rate": 10, "user_classes": [MarketMSRequests]},
        {"duration": 435, "users": 10, "spawn_rate": 5, "user_classes": [UserMSRequests, GachaMSRequests]},
    ]

    def tick(self):
        run_time = self.get_run_time()

        for stage in self.stages:
            if run_time < stage["duration"]:
                try:
                    tick_data = (stage["users"], stage["spawn_rate"], stage["user_classes"])
                except:
                    tick_data = (stage["users"], stage["spawn_rate"])
                return tick_data
        return None

# ==============================================================================
# Utils
# ==============================================================================

def generate_random_string(length=28):
    """
    Generate a cryptographically secure random string

    Args:
        length (int): Desired length of the random string (default is 28)

    Returns:
        str: Hexadecimal random string
    """
    # Generate random bytes
    random_bytes = secrets.token_bytes(length // 2)

    # Convert bytes to hex string
    hex_string = random_bytes.hex()

    return hex_string

def parse_jwt(token: str, algorithms=["H512"]) -> dict | None:
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

# ==============================================================================
# OAuth 2.0 Utils
# ==============================================================================

def sha256(plain):
    """
    Create a SHA-256 hash of the input string
    Returns bytes of the hash
    """
    return hashlib.sha256(plain.encode('utf-8')).digest()

def base64_urlencode(data):
    """
    Convert bytes to URL-safe base64 encoding
    Removes padding and replaces standard base64 characters
    """
    # Base64 encode and convert to string
    base64_encoded = base64.b64encode(data).decode('utf-8')

    # URL-safe modifications
    return (base64_encoded
        .replace('+', '-')  # Replace + with -
        .replace('/', '_')  # Replace / with _
        .replace('=', '')   # Remove padding
    )

def pkce_challenge_from_verifier(v):
    """
    Generate PKCE code challenge from code verifier
    """
    hashed = sha256(v)
    return base64_urlencode(hashed)

def generate_code_verifier_and_chall():
    verifier = generate_random_string(50)
    return verifier, pkce_challenge_from_verifier(verifier)
