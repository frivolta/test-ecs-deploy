version: "3.9"
services:
postgres:
image: postgres:12-alpine
environment:
- POSTGRES_USER=root
- POSTGRES_PASSWORD=secret
- POSTGRES_DB=apecalendar
ports:
- "5433:5432"
birdie:
build:
context: ./birdie
dockerfile: Dockerfile.dev
ports:
- "8080:8080"
depends_on:
- postgres
entrypoint: [ "/app/wait-for.sh", "postgres:5432", "--", "/app/start.sh" ]
command: [ "air", "run", "main.go" ]
client:
build:
context: ./client
dockerfile: Dockerfile
stdin_open: true
depends_on:
- birdie
ports:
- '3002:3000'
environment:
- REACT_APP_API_KEY="AIzaSyCqHIEOk0CQk3lBkAA4ldfwNFzpvoeJ9Pg"
- REACT_APP_AUTH_DOMAIN="apecalendar-b4079.firebaseapp.com"
- REACT_APP_PROJECT_ID="apecalendar-b4079"
- REACT_APP_STORAGE_BUCKET="apecalendar-b4079.appspot.com"
- REACT_APP_MESSAGING_SENDER_ID="680474127642"
- REACT_APP_APP_ID="1:680474127642:web:58eb50aaa397f376b510a3"
- REACT_APP_MEASUREMENT_ID="G-023CYZY9FT"
- REACT_APP_DOMAIN="dev-iv8n6772.us.auth0.com"
- REACT_APP_AUDIENCE="https://apecalendar-dev.com"
- REACT_APP_CLIENT_ID="g2qwqnmFAmXkfsJYnvXSQuvV7X5bc4fu"
- SKIP_PREFLIGHT_CHECK=true
volumes:
- ./client/src:/app/src
- /app/node_modules