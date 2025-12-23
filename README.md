# Go Telegram System Monitor

A robust, concurrent system monitoring service written in Go. This application runs in the background, sending daily system resource reports to a Telegram chat and providing real-time alerts if memory or disk usage exceeds critical thresholds.

## 🚀 Features

*   **Daily Reports:** Automatically sends a summary of CPU, RAM, and Disk usage at **07:30** and **17:30** every day.
*   **Real-time Alerting:** Continuously checks system resources (every 5 minutes).
*   **Threshold Warnings:** Sends an immediate alert if Memory or Disk usage exceeds **80%**.
*   **Concurrent Design:** Uses Go Goroutines to handle scheduling and monitoring independently.
*   **Systemd Integration:** Ready for production deployment as a Linux background service.

## 🛠️ Prerequisites

*   **Go** 1.24 or higher
*   **Telegram Bot Token** (Get one from [@BotFather](https://t.me/BotFather))
*   **Telegram Chat ID** (Get yours from [@userinfobot](https://t.me/userinfobot))

## ⚙️ Installation & Setup

1.  **Clone the repository**
    ```bash
    git clone <your-repo-url>
    cd go-telegram-monitor
    ```

2.  **Configure Environment Variables**
    Copy the example file and edit it with your credentials:
    ```bash
    cp .env.example .env
    nano .env
    ```
    *   `TELEGRAM_BOT_TOKEN`: Your bot token.
    *   `TELEGRAM_CHAT_ID`: Your personal or group chat ID.

3.  **Install Dependencies**
    ```bash
    go mod tidy
    ```

4.  **Build the Binary**
    Compile the application into a standalone executable named `monitor`:
    ```bash
    go build -o monitor .
    ```

## 📦 Deployment (Systemd)

To run this application as a background service on Linux (Ubuntu/Debian/CentOS):

1.  **Edit the Service File**
    Open `go-telegram-monitor.service` and ensure the paths match your installation directory and username.
    ```ini
    User=your_username
    WorkingDirectory=/path/to/go-telegram-monitor
    ExecStart=/path/to/go-telegram-monitor/monitor
    ```

2.  **Install the Service**
    ```bash
    sudo cp go-telegram-monitor.service /etc/systemd/system/
    sudo systemctl daemon-reload
    ```

3.  **Enable and Start**
    ```bash
    sudo systemctl enable go-telegram-monitor
    sudo systemctl start go-telegram-monitor
    ```

4.  **Check Status & Logs**
    *   Check status: `sudo systemctl status go-telegram-monitor`
    *   View logs: `journalctl -u go-telegram-monitor -f`
    *   Stop service: `sudo systemctl stop go-telegram-monitor`

## 🏗️ Project Structure

*   `main.go`: Application entry point, scheduler setup, and monitoring logic.
*   `monitor`: Compiled executable binary (ignored by Git).
*   `go-telegram-monitor.service`: Systemd configuration file.
*   `.env`: Configuration file for secrets (ignored by Git).

## 👨‍💻 Tech Stack

*   **Language:** Go (Golang)
*   **Libraries:**
    *   `gopsutil`: System metrics collection.
    *   `gocron`: Job scheduling.
    *   `godotenv`: Environment variable management.

---
*Created as a portfolio project to demonstrate Go concurrency and system automation patterns.*
