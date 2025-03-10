#!/bin/bash

# Create directory for certificates
mkdir -p certs

# Generate root certificate
openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 \
  -keyout certs/rootCA.key -out certs/rootCA.pem \
  -subj "/C=EN/CN=Local-Root-CA"

# Create configuration file
cat > certs/domains.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = pets.local
EOF

# Generate key and CSR for domain
openssl req -new -nodes -newkey rsa:2048 \
  -keyout certs/pets.local.key -out certs/pets.local.csr \
  -subj "/C=EN/ST=State/L=City/O=Organization/CN=pets.local"

# Create certificate for domain, signed by root CA
openssl x509 -req -sha256 -days 1024 \
  -in certs/pets.local.csr \
  -CA certs/rootCA.pem -CAkey certs/rootCA.key -CAcreateserial \
  -extfile certs/domains.ext \
  -out certs/pets.local.crt

echo "Certificates have been generated. Now add the rootCA.pem to your browser's trusted certificates."
