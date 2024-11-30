#!/usr/bin/env bash

set -eo pipefail
REGENERATE="${1:-false}"
set -u

if [[ ! -f "./cert.pem" || ! -f "./key.pem" ]]; then
  echo "CA certificate or key not found. Exiting."
  exit 1
fi

declare -Ar psg_users=(
  ["user/postgres"]="user-db"
  ["admin/postgres"]="admin-db"
  ["gacha/postgres"]="gacha-db"
  ["market/postgres"]="market-db"
)

declare -Ar rds_users=(
  ["auth/redis"]="sessions-redis"
  ["market/redis"]="market-timed-events-redis"
)

u1="$(declare -p psg_users)"
u2="$(declare -p rds_users)"
users_as_string="${u1:0:${#u1}-1} ${u2:23}"
eval "declare -A users="${users_as_string#*=}

cd ../

if [[ "$REGENERATE" == "-r" ]]; then
  echo "Removing existing certificates..."
  for user in "${!users[@]}"; do
    BASE_DIR="$user/certs"
    if [[ -f "$BASE_DIR/server.crt" || -f "$BASE_DIR/server.key" ]]; then
      sudo rm "$BASE_DIR"/server.* || { echo "Failed to remove existing certificate files for $user"; exit 1; }
    fi
  done
fi

for user in "${!users[@]}"; do
  CN="${users[$user]}"

  BASE_DIR="$user/certs"
  mkdir -p "$BASE_DIR"

  openssl genrsa -out "$BASE_DIR/server.key" 4096 || { echo "Failed to generate key for $user"; exit 1; }

  openssl req -new -key "$BASE_DIR/server.key" -out "$BASE_DIR/server.csr" -subj "/CN=$CN" || { echo "Failed to generate CSR for $user"; exit 1; }

  openssl x509 -req -extfile <(printf "subjectAltName=DNS:$CN") -in "$BASE_DIR/server.csr" -CA ./cacerts/cert.pem -CAkey ./cacerts/key.pem -CAcreateserial -out "$BASE_DIR/server.crt" -days 365 || { echo "Failed to generate certificate for $user"; exit 1; }

  rm "$BASE_DIR/server.csr" || { echo "Failed to remove CSR file for $user"; exit 1; }
done

rm ./cacerts/cert.srl

read -p "Do you want to change ownership of the certificates [y/n]? " answer

case $answer in
  [Yy]* )
    for user in "${!psg_users[@]}"; do
      BASE_DIR="$user/certs"
      sudo chown 999:999 "$BASE_DIR/server."* || { echo "Failed to change ownership for $user"; exit 1; }
      echo "Ownership changed for $user."
    done
    for user in "${!rds_users[@]}"; do
      BASE_DIR="$user/certs"
      sudo chown 999:1000 "$BASE_DIR/server."* || { echo "Failed to change ownership for $user"; exit 1; }
      echo "Ownership changed for $user."
    done
    ;;
  * )
    echo "Ownership not changed for any user."
    ;;
esac

echo "All certificates generated successfully."
