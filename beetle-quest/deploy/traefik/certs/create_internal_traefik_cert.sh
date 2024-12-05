#!/bin/sh

openssl genpkey -algorithm RSA -out client_internal_key.pem -pkeyopt rsa_keygen_bits:2048
openssl req -new -key client_internal_key.pem -out client_internal_cert.pem -subj "/C=IT/ST=State/L=City/O=Organization/OU=Department/CN=reverse-proxy"
openssl x509 -req -in client_internal_cert.pem -CA ../../cacerts/cert.pem -CAkey ../../cacerts/key.pem -CAcreateserial -out client_internal_cert.pem -days 365 -extfile <(printf "subjectAltName=DNS:admin-reverse-proxy,DNS:reverse-proxy")
