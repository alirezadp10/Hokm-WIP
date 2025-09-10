# Hokm 🎴

A multiplayer implementation of the Persian card game **Hokm**, built with a hybrid tech stack.  
The project combines **Go** for backend logic, **JavaScript/HTML/CSS** for the web interface, and **Lua** for game state definitions.  
It is containerized with Docker and supports environment-based configuration for easy deployment.

---

## ✨ Features

- **Multiplayer Hokm Gameplay** – Card shuffling, turn management, and game rules implemented.  
- **Hybrid Tech Stack** – Go server, JavaScript frontend, HTML templates, and Lua scripting.  
- **Template-Driven UI** – HTML/CSS templates for rendering the game state.  
- **Configurable Environment** – `.env` file for secrets and runtime settings.  
- **Containerized Deployment** – Run with Docker and Docker Compose.  
- **Persistent Storage Ready** – Designed to integrate with databases or caches (via `.env`).

---

## 📂 Project Structure

```
hokm/
├── assets/                 # Static assets (images, CSS, JS)
├── cmd/
│   └── server/             # Go server entrypoint
├── internal/               # Internal Go packages
├── pkg/                    # Shared Go packages
├── templates/              # HTML templates for UI
├── game-states.lua         # Lua-based game state logic
├── .env.example            # Example environment variables
├── Dockerfile              # Docker build instructions
├── docker-compose.yaml     # Compose file for services
├── go.mod / go.sum         # Go module definitions
└── note                    # Project notes
```

---

## ⚙️ Getting Started

### Prerequisites

- [Go 1.18+](https://go.dev/)  
- [Node.js](https://nodejs.org/) (if building frontend assets)  
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/) (optional)  

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/alirezadp10/hokm.git
   cd hokm
   ```

2. Copy the environment file:

   ```bash
   cp .env.example .env
   ```

   Update `.env` with your configuration (server port, database URL, Redis, etc.).

3. Install Go dependencies:

   ```bash
   go mod download
   ```

---

### Running Locally

Start the Go server:

```bash
go run cmd/server/main.go
```

By default, the app runs on `http://localhost:8080`.

---

### Running with Docker

Build and run the containerized setup:

```bash
docker-compose up --build
```

This starts the Hokm server and any configured dependencies.

---

## 🕹️ Usage

- Open the application in your browser (default: `http://localhost:8080`).  
- Create or join a Hokm match.  
- Play turns in real time—the server manages game state and enforces rules.  

---

## 🛠️ Tech Stack

- **Go** – Backend server and game management  
- **JavaScript / HTML / CSS** – Web frontend  
- **Lua** – Game state definitions and logic scripting  
- **Docker** – Containerization for deployment  
- *(Optional)* Redis / PostgreSQL if you integrate persistence  

---

## 🤝 Contributing

Contributions are welcome! To get started:

1. Fork the repo  
2. Create a new branch (`feature/my-feature`)  
3. Commit your changes  
4. Submit a Pull Request  

---

## 📜 License

This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.
