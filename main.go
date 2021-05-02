package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"gopkg.in/natefinch/lumberjack.v2"
)

type WebhookEvent struct {
	ID        string
	Type      string
	Signature string
	Algorithm string
	Payload   []byte
}

func NewWebhookEvent(r *http.Request) *WebhookEvent {
	w := new(WebhookEvent)
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var pld []byte
	var obj map[string]interface{}
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		s := string(bytes)
		// strip all newlines so it logs correctly
		n := strings.Replace(s, "\n", " ", -1)
		pld = []byte(n)
	} else {
		pld, err = json.Marshal(obj)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(pld) > 0 {
		w.Payload = pld
	} else {
		w.Payload = []byte("-")
	}

	if i := r.Header.Get("X-Nexus-Webhook-Delivery"); i != "" {
		w.ID = i
	} else {
		w.ID = "-"
	}
	if i := r.Header.Get("X-Nexus-Webhook-ID"); i != "" {
		w.Type = i
	} else {
		w.Type = "-"
	}
	if i := r.Header.Get("X-Nexus-Webhook-Signature"); i != "" {
		w.Signature = i
	} else {
		w.Signature = "-"
	}
	if i := r.Header.Get("X-Nexus-Webhook-Signature-Algorithm"); i != "" {
		w.Algorithm = i
	} else {
		w.Algorithm = "-"
	}

	return w
}

func isWritable(path string) bool {
	return syscall.Access(path, 2) == nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		we := NewWebhookEvent(r)
		msg := fmt.Sprintf(
			"%s %s %s %s %s",
			we.ID, we.Type,
			we.Signature,
			we.Algorithm,
			we.Payload,
		)
		log.Print(msg)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var (
	logger *lumberjack.Logger

	// cmd flags
	portFlag int
	logPath  string
)

func init() {
	fp, _ := os.Getwd()
	fn := "events.log"
	defaultPath := fmt.Sprintf("%s/%s", fp, fn)
	defaultPort := 3000

	flag.IntVar(&portFlag, "port", defaultPort, "listen port")
	flag.StringVar(&logPath, "path", defaultPath, "log file path")
	flag.Parse()

	pathOnly := fmt.Sprint(filepath.Dir(logPath))
	if !isWritable(pathOnly) {
		log.Fatalf("Fatal: Not authorised to write to %s", logPath)
	}

	logger = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    128,  // megabytes
		MaxBackups: 0,    //
		MaxAge:     7,    // days
		Compress:   true, // disabled by default
	}
	log.SetOutput(logger)
}

func main() {
	http.HandleFunc("/callback", callbackHandler)

	fmt.Printf("Logging events to %s\n", logPath)
	fmt.Printf("Listening for events at http://localhost:%d\n", portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portFlag), nil))
}
