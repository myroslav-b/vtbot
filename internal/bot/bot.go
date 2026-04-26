package bot

import (
	"fmt"
	"io"
	"log"
	"time"

	"vtbot/internal/config"
	"vtbot/internal/utils"
	"vtbot/internal/virustotal"

	tele "gopkg.in/telebot.v3"
)

type Job struct {
	Context   tele.Context
	Document  *tele.Document
	StatusMsg *tele.Message
}

type Bot struct {
	bot *tele.Bot
	vt  *virustotal.Client
	cfg *config.Config
}

func New(cfg *config.Config, vt *virustotal.Client) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot: b,
		vt:  vt,
		cfg: cfg,
	}, nil
}

func (b *Bot) Start() {
	jobQueue := make(chan Job, 100)

	go b.worker(jobQueue)

	b.bot.Handle(tele.OnDocument, func(c tele.Context) error {
		doc := c.Message().Document

		if utils.ShouldIgnoreFile(doc.FileName) {
			return nil
		}

		if doc.FileSize > b.cfg.MaxFileSize {
			return c.Reply(fmt.Sprintf("⚠️ Файл %s завеликий (>%dMB). Пропускаємо.", doc.FileName, b.cfg.MaxFileSize/(1024*1024)))
		}

		statusMsg, err := b.bot.Reply(c.Message(), fmt.Sprintf("⏳ Файл %s додано в чергу на перевірку...", doc.FileName))
		if err != nil {
			log.Println("Не вдалося відправити повідомлення про чергу:", err)
			return err
		}

		select {
		case jobQueue <- Job{Context: c, Document: doc, StatusMsg: statusMsg}:
			return nil
		case <-time.After(2 * time.Second):
			b.bot.Delete(statusMsg)
			return c.Reply("⚠️ Черга переповнена. Спробуйте пізніше.")
		}
	})

	log.Println("Security Bot запущено...")
	b.bot.Start()
}

func (b *Bot) worker(jobs <-chan Job) {
	ticker := time.NewTicker(b.cfg.RequestInterval)
	defer ticker.Stop()

	for job := range jobs {
		<-ticker.C
		b.processJob(job)
	}
}

func (b *Bot) processJob(job Job) {
	doc := job.Document
	statusMsg := job.StatusMsg

	_, err := b.bot.Edit(statusMsg, fmt.Sprintf("🔍 Аналізую %s...", doc.FileName))
	if err != nil {
		log.Println("Не вдалося оновити повідомлення:", err)
	}

	fileReader, err := b.bot.File(&doc.File)
	if err != nil {
		b.bot.Edit(statusMsg, "❌ Помилка завантаження файлу з Telegram.")
		return
	}
	defer fileReader.Close()

	content, err := io.ReadAll(fileReader)
	if err != nil {
		b.bot.Edit(statusMsg, "❌ Помилка читання файлу.")
		return
	}

	hash := utils.CalculateSHA256(content)

	report, found, err := b.vt.GetReportByHash(hash)
	if err != nil {
		log.Printf("API Error (Hash Check): %v", err)
		b.bot.Edit(statusMsg, "❌ Помилка з'єднання з VirusTotal.")
		return
	}

	if found {
		b.sendReport(statusMsg, doc.FileName, report, "✅ Знайдено в базі (без завантаження)")
		return
	}

	b.bot.Edit(statusMsg, "📤 Файл невідомий. Завантажую на VirusTotal...")

	analysisID, err := b.vt.UploadFile(doc.FileName, content)
	if err != nil {
		log.Printf("Upload error for %s: %v", doc.FileName, err)
		b.bot.Edit(statusMsg, "❌ Помилка завантаження файлу на VirusTotal.")
		return
	}

	// Polling виконується в окремій горутині, щоб не блокувати worker.
	// PollAnalysis має свій внутрішній rate limiting (15s між запитами).
	go func() {
		finalReport, err := b.vt.PollAnalysis(analysisID)
		if err != nil {
			log.Printf("Poll error for %s (analysis %s): %v", doc.FileName, analysisID, err)
			b.bot.Edit(statusMsg, "⚠️ Час очікування аналізу вичерпано або помилка.")
			return
		}

		b.sendReport(statusMsg, doc.FileName, finalReport, "🆕 Новий аналіз")
	}()
}

func (b *Bot) sendReport(msg *tele.Message, filename string, report *virustotal.VTResponse, note string) {
	safeName := utils.EscapeMarkdown(filename)
	stats := report.Data.Attributes.Stats

	// Якщо stats порожній (сума 0), можливо це LastAnalysisStats (з файлу)
	if stats.Malicious == 0 && stats.Suspicious == 0 && stats.Harmless == 0 && stats.Undetected == 0 {
		stats = report.Data.Attributes.LastAnalysisStats
	}

	emoji := "✅"
	verdict := "Чистий"

	if stats.Malicious > 0 {
		emoji = "⛔️"
		verdict = fmt.Sprintf("ЗАГРОЗА (%d детектів)", stats.Malicious)
	} else if stats.Suspicious > 0 {
		emoji = "⚠️"
		verdict = "Підозрілий"
	}

	var text string
	if b.cfg.VerboseOutput {
		text = fmt.Sprintf(
			"%s **Результат:** %s\n"+
				"📂 Файл: `%s`\n"+
				"📊 Інфо: %s\n\n"+
				"🔴 Malicious: %d\n"+
				"🟠 Suspicious: %d\n"+
				"🟢 Harmless/Undetected: %d\n",
			emoji, verdict, safeName, note,
			stats.Malicious, stats.Suspicious, stats.Harmless+stats.Undetected,
		)
	} else {
		text = fmt.Sprintf(
			"📂 Файл: `%s`\n"+
				"%s Результат: %s\n",
			safeName, emoji, verdict,
		)
	}

	b.bot.Edit(msg, text, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}
