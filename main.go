package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
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

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	log.Println("Scheduler created.")

	_, err = s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(7, 30, 0),
				gocron.NewAtTime(17, 30, 0),
			),
		),
		gocron.NewTask(runScheduledReport, telegramBotToken, telegramChatID),
	)
	if err != nil {
		log.Fatalf("Failed to schedule daily report job: %v", err)
	}
	log.Println("Scheduled daily reports for 07:30 and 17:30.")

	go runContinuousMonitoring(telegramBotToken, telegramChatID)
	log.Println("Started continuous monitoring for high-resource alerts.")

	s.Start()
	log.Println("Scheduler started.")
	select {}
}

func runScheduledReport(botToken, chatID string) {
	log.Println("Running scheduled report...")
	metrics, err := getSystemMetrics()
	if err != nil {
		log.Printf("Scheduled report: Failed to get system metrics: %v", err)
		return
	}

	reportMessage := "📊 Daily System Metrics Report 📊\n\n" + metrics
	err = sendTelegramMessage(botToken, chatID, reportMessage)
	if err != nil {
		log.Printf("Scheduled report: Failed to send Telegram message: %v", err)
	}
}

func runContinuousMonitoring(botToken, chatID string) {
	ticker := time.NewTicker(5 * time.Minute)
	log.Println("Continuous monitoring will check every 5 minutes.")

	for range ticker.C {
		log.Println("Running continuous monitoring check...")

		vm, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Monitor: Failed to get memory usage: %v", err)
			continue
		}

		du, err := disk.Usage("/")
		if err != nil {
			log.Printf("Monitor: Failed to get disk usage: %v", err)
			continue
		}

		alertMessage := ""
		if vm.UsedPercent > 80.0 {
			alertMessage += fmt.Sprintf("🚨 High Memory Usage Alert: %.2f%%\n", vm.UsedPercent)
		}
		if du.UsedPercent > 80.0 {
			alertMessage += fmt.Sprintf("🚨 High Disk Usage Alert: %.2f%%", du.UsedPercent)
		}

		if alertMessage != "" {
			log.Println("High resource usage detected, sending alert.")
			err := sendTelegramMessage(botToken, chatID, alertMessage)
			if err != nil {
				log.Printf("Monitor: Failed to send alert message: %v", err)
			}
		} else {
			log.Println("Monitoring check complete, all systems nominal.")
		}
	}
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
		"Disk Usage:\n  Total: %d GB\n  Used: %d GB\n  Used Percent: %.2f%%",
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

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK status from Telegram API: %s", resp.Status)
	}

	log.Printf("Successfully sent message to chat ID %s.", chatID)
	return nil
}
