#!/usr/bin/env bash
# JOSE header and JWT payload
VAULT_ADDR="http://localhost:8200"
HEADER=$(echo  '{"alg": "ES256","typ": "JWT", "pad": "....."}' | jq -c '.' | base64)
PAYLOAD=$(echo '{"sub": "1234567890", "name": "John Doe", "email": "john.doe@google.com"}' | jq -c '.' | base64)

# Create a key in Vault.
vault secrets enable transit >/dev/null 2>&1
vault write transit/keys/mykey exportable=true type=ecdsa-p256 >/dev/null 2>&1

# Prepare header and payload for signing
MESSAGE=$(echo -n "$HEADER.$PAYLOAD" | openssl base64 -A)

# Sign the message using JWS marshaling type, and remove the vault key prefix
JWS=$(vault write -format=json transit/sign/mykey input=$MESSAGE marshaling_algorithm=jws | jq -r .data.signature | cut -d ":" -f3)

# Combine to build the JWT
JWT="$HEADER.$PAYLOAD.$JWS"

# Export the the key and print out the public key portion
vault read -format=json transit/export/signing-key/mykey/1 | jq -r '.data.keys."1"' > /tmp/privkey
KEY=$(openssl ec -in /tmp/privkey -pubout 2>/dev/null | base64)

echo "{\"jwt\": \"${JWT}\", \"pub\": \"${KEY}\"}"


HEADER=$(echo  '{"alg": "ES256","typ": "JWT", "pad": "....."}' | jq -c '.' | base64)
PAYLOAD=$(echo '{"sub": "1234567890", "name": "Alice Dee", "email": "alice.dee@apple.com"}' | jq -c '.' | base64)
MESSAGE=$(echo -n "$HEADER.$PAYLOAD" | openssl base64 -A)
JWS=$(vault write -format=json transit/sign/mykey input=$MESSAGE marshaling_algorithm=jws | jq -r .data.signature | cut -d ":" -f3)
JWT="$HEADER.$PAYLOAD.$JWS"
echo "{\"jwt\": \"${JWT}\", \"pub\": \"${KEY}\"}"
# You should be able to successfully decode the JWT on https://jwt.io
