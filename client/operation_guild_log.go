package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ao-data/albiondata-client/log"
)

// --- AYARLAR ---
const DiscordWebhookURL = "BURAYA_WEBHOOK_LINKINIZI_YAZIN"

// Cursor Dosyası
const CursorFile = "last_log_cursor.txt"

// Bildirim Eşiği (Örn: 10 Milyon Silver)
// Şu an test için 10.000.000 olarak ayarlı.
const NotificationThreshold = 10000000

// --- CONSTANTS ---
const ActionTypeWithdraw = 2 // Para Çekme İşlem Kodu (Tahmini)

type OperationGuildLogResponse struct {
	PlayerNames []string `mapstructure:"0"`
	ActionTypes []int    `mapstructure:"1"`
	Amounts     []int64  `mapstructure:"3"`
	Timestamps  []int64  `mapstructure:"4"`
}

type DiscordMessage struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func (op OperationGuildLogResponse) Process(state *albionState) {
	lastProcessedTime := loadLastCursor()
	newMaxTime := lastProcessedTime

	log.Infof(">>> Loglar taranıyor... (Son Kayıt Ticks: %d)", lastProcessedTime)

	for i := 0; i < len(op.PlayerNames); i++ {
		// Dizi sınırlarını kontrol et
		if i >= len(op.Amounts) || i >= len(op.ActionTypes) || i >= len(op.Timestamps) {
			break
		}

		currentTicks := op.Timestamps[i]

		// Eskileri atla
		if currentTicks <= lastProcessedTime {
			continue
		}

		if currentTicks > newMaxTime {
			newMaxTime = currentTicks
		}

		// --- VERİ DÜZENLEME ---
		playerName := op.PlayerNames[i]
		rawAmount := op.Amounts[i]
		actionType := op.ActionTypes[i]

		// 10000 birim = 1 Silver
		realAmount := rawAmount / 10000

		// --- FİLTRELEME (BURASI DEĞİŞTİ) ---
		// 1. İşlem Tipi 2 mi? (Para Çekme)
		// 2. Miktar Eşik değerinden büyük mü?
		if actionType == ActionTypeWithdraw && realAmount >= NotificationThreshold {

			log.Infof("⚠️ ŞÜPHELİ İŞLEM: %s - %s Silver", playerName, formatWithDots(realAmount))

			// Zamanı hesapla
			logTime := ticksToTime(currentTicks)

			// Discord'a gönder
			go sendDiscordAlert(playerName, realAmount, logTime)
		}
	}

	// Sadece yeni bir işlem gördüysek cursor'u güncelle
	if newMaxTime > lastProcessedTime {
		saveCursor(newMaxTime)
	}
}

// --- YARDIMCI FONKSİYONLAR ---

// C# Ticks (0001 yılı) -> Go Time (1970 yılı)
func ticksToTime(ticks int64) time.Time {
	const ticksAtEpoch = 621355968000000000
	const ticksPerSecond = 10000000

	unixTicks := ticks - ticksAtEpoch

	seconds := unixTicks / ticksPerSecond
	nanos := (unixTicks % ticksPerSecond) * 100

	return time.Unix(seconds, nanos)
}

func formatWithDots(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits--
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}

func loadLastCursor() int64 {
	data, err := ioutil.ReadFile(CursorFile)
	if err != nil {
		return 0
	}
	t, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0
	}
	return t
}

func saveCursor(t int64) {
	_ = ioutil.WriteFile(CursorFile, []byte(fmt.Sprintf("%d", t)), 0644)
}

func sendDiscordAlert(player string, amount int64, logTime time.Time) {
	formattedAmount := formatWithDots(amount)

	// Tarihi TR formatında göster
	localTime := logTime.Local().Format("02.01.2006 15:04:05")

	message := fmt.Sprintf(
		"🚨 **ALARM: BANKADAN YÜKSEK MİKTAR ÇEKİLDİ (10M+)** 🚨\n\n"+
			"👤 **Oyuncu:** `%s`\n"+
			"💸 **Çekilen Miktar:** `%s Silver`\n"+
			"📅 **İşlem Zamanı(UTC+3):** `%s`\n"+
			"⚠️ **Aksiyon:** Lütfen bu çekimin sebebini oyuncuya sorunuz.\n"+
			"--------------------------------------",
		player, formattedAmount, localTime,
	)

	payload := DiscordMessage{
		Content:  message,
		Username: "Guild Bank",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", DiscordWebhookURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	client.Do(req)
}
