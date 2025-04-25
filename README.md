# How to Run the Program

## Requirements
- Install [Docker](https://www.docker.com/) (if not already installed)
- Install [ngrok](https://ngrok.com/)
- Create a [Pusher](https://pusher.com/) account and get your credentials

## Setup Instructions

### 1. Environment Variables
- Create a ` .env` file using ` .example.env` as a template.
- Update using the credential from you Pusher account
- Do this for both:
  - the main program at [` .example.env`](./.example.env)
  - the frontend located at [`client/my-qr-client/.example.env`](./client/my-qr-client/.example.env)

### 2. Start Redis with Docker
Make sure Redis is running on the default port (6379). You can use the following Docker command:
```bash
docker run --name my-redis -p 6379:6379 -d redis
```

### 3. Configure ngrok
Update your `ngrok.yml` file (example Windows path):
```
C:\Users\<your-username>\.ngrok2\ngrok.yml
```

Example config:
```yaml
version: "2"
authtoken: YOUR_NGROK_TOKEN  # optional if already authenticated

tunnels:
  frontend:
    addr: 5173
    proto: http
  backend:
    addr: 1323
    proto: http
```

Then run:
```bash
ngrok start --all
```

### 4. Update URLs
After ngrok starts, copy the generated links for both frontend and backend and update:

- In [`auth/auth.go`](./auth/auth.go):
```go
var baseURL = "https://<your-ngrok-backend-url>" // CHANGE THIS
```

- In [`client/my-qr-client/main.js`](./client/my-qr-client/main.js):
```js
const baseURL = "https://<your-ngrok-backend-url>" // CHANGE THIS
```

### 5. Run the Applications

#### Backend
```bash
go run .
```

#### Frontend
```bash
cd client/my-qr-client
npm run dev
```