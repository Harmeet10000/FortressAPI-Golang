# Generate a self-signed certificate and private key
openssl req -x509 -newkey rsa:2048 -nodes -keyout certs/key.pem -out certs/cert.pem -days 365 -config certs/openssl.cnf
