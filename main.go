package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in environment")
	}

	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramChatID == "" {
		log.Fatal("TELEGRAM_CHAT_ID is not set in environment")
	}

	http.HandleFunc("/send-metrics", func(w http.ResponseWriter, r *http.Request) {
		handleSendMetrics(w, r, telegramBotToken, telegramChatID)
	})

	fmt.Println("Server starting on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func handleSendMetrics(w http.ResponseWriter, r *http.Request, telegramBotToken, telegramChatID string) {
	fmt.Println("Received request to /send-metrics")

	metrics, err := getSystemMetrics()
	if err != nil {
		log.Printf("Error getting system metrics: %v", err)
		http.Error(w, "Failed to retrieve system metrics", http.StatusInternalServerError)
		return
	}

	fmt.Println("--- System Metrics ---")
	fmt.Println(metrics)
	fmt.Println("----------------------")

	err = sendTelegramMessage(telegramBotToken, telegramChatID, metrics)
	if err != nil {
		log.Printf("Error sending Telegram message: %v", err)
		http.Error(w, "Failed to send metrics to Telegram", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Metrics collected and sent to Telegram successfully!"))
	log.Println("Metrics collected and sent to Telegram successfully!")
}

func getSystemMetrics() (string, error) {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", fmt.Errorf("failed to get cpu usage: %w", err)
	}
	cpuUsageStr := fmt.Sprintf("CPU Usage: %.2f%%\n", cpuUsage[0])

	vm, err := mem.VirtualMemory()
	if err != nil {
		return "", fmt.Errorf("failed to get memory usage: %w", err)
	}
	memUsageStr := fmt.Sprintf(
		"Memory Usage:\n  Total: %d MB\n  Used: %d MB\n  Used Percent: %.2f%%\n",
		vm.Total/(1024*1024),
		vm.Used/(1024*1024),
		vm.UsedPercent,
	)

	du, err := disk.Usage("/")
	if err != nil {
		return "", fmt.Errorf("failed to get disk usage: %w", err)
	}
	diskUsageStr := fmt.Sprintf(
		"Disk Usage:\n  Total: %d GB\n  Used: %d GB\n  Used Percent: %.2f%%\n",
		du.Total/(1024*1024*1024),
		du.Used/(1024*1024*1024),
		du.UsedPercent,
	)

	return cpuUsageStr + memUsageStr + diskUsageStr, nil
}

func sendTelegramMessage(botToken, chatID, message string) error {
	telegramAPIURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", telegramAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create Telegram request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK status from Telegram API: %s", resp.Status)
	}

	log.Println("Message sent to Telegram successfully.")
	return nil
}
