package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/Rhymen/go-whatsapp"
)

func main() {
    // Membuat koneksi
    wac, err := whatsapp.NewConn(20 * time.Second)
    if err != nil {
        log.Fatalf("Error creating connection: %v", err)
    }

    var sess whatsapp.Session
    sessionFile := "./session/session.json"

    // Memulihkan sesi jika ada file sesi
    if _, err := os.Stat(sessionFile); err == nil {
        file, err := os.Open(sessionFile)
        if err != nil {
            log.Fatalf("Error opening session file: %v", err)
        }
        defer file.Close()

        decoder := json.NewDecoder(file)
        err = decoder.Decode(&sess)
        if err != nil {
            log.Fatalf("Error decoding session file: %v", err)
        }

        // Memulihkan sesi
        newSess, err := wac.RestoreWithSession(sess)
        if err != nil {
            log.Fatalf("Error restoring session: %v", err)
        }
        sess = newSess
    } else {
        // Login jika tidak ada file sesi
        qrChan := make(chan string)
        go func() {
            qr := <-qrChan
            fmt.Printf("QR Code: %s\n", qr)
            // Tampilkan QR code atau simpan untuk dipindai
        }()

        sess, err = wac.Login(qrChan)
        if err != nil {
            log.Fatalf("Login error: %v", err)
        }

        // Simpan sesi ke file
        file, err := os.Create(sessionFile)
        if err != nil {
            log.Fatalf("Error creating session file: %v", err)
        }
        defer file.Close()

        encoder := json.NewEncoder(file)
        err = encoder.Encode(sess)
        if err != nil {
            log.Fatalf("Error encoding session file: %v", err)
        }
    }

    // Menambahkan handler untuk pesan masuk
    wac.AddHandler(myHandler{})

    // Menunggu secara tak terbatas
    select {}
}

type myHandler struct{}

func (myHandler) HandleError(err error) {
    fmt.Fprintf(os.Stderr, "%v", err)
}

func (myHandler) HandleTextMessage(message whatsapp.TextMessage) {
    fmt.Println("Received Text Message:", message.Text)
}

func (myHandler) HandleImageMessage(message whatsapp.ImageMessage) {
    fmt.Println("Received Image Message")
}

func (myHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
    fmt.Println("Received Document Message")
}

func (myHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
    fmt.Println("Received Video Message")
}

func (myHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
    fmt.Println("Received Audio Message")
}

func (myHandler) HandleJsonMessage(message string) {
    fmt.Println("Received JSON Message:", message)
}

func (myHandler) HandleContactMessage(message whatsapp.ContactMessage) {
    fmt.Println("Received Contact Message")
}

func (myHandler) HandleBatteryMessage(message whatsapp.BatteryMessage) {
    fmt.Println("Received Battery Message")
}

func (myHandler) HandleNewContact(contact whatsapp.Contact) {
    fmt.Println("Received New Contact:", contact)
}

// Mengirim pesan teks
func sendTextMessage(wac *whatsapp.Conn) {
    text := whatsapp.TextMessage{
        Info: whatsapp.MessageInfo{
            RemoteJid: "0123456789@s.whatsapp.net",
        },
        Text: "Hello Whatsapp",
    }

    id, err := wac.Send(text)
    if err != nil {
        log.Fatalf("Error sending text message: %v", err)
    }
    fmt.Printf("Message sent with ID: %s\n", id)
}

// Mengirim pesan kontak
func sendContactMessage(wac *whatsapp.Conn) {
    contactMessage := whatsapp.ContactMessage{
        Info: whatsapp.MessageInfo{
            RemoteJid: "0123456789@s.whatsapp.net",
        },
        DisplayName: "Luke Skylwalker",
        Vcard: "BEGIN:VCARD\nVERSION:3.0\nN:Skyllwalker;Luke;;\nFN:Luke Skywallker\nitem1.TEL;waid=0123456789:+1 23 456789789\nitem1.X-ABLabel:Mobile\nEND:VCARD",
    }

    id, err := wac.Send(contactMessage)
    if err != nil {
        log.Fatalf("Error sending contact message: %v", err)
    }
    fmt.Printf("Contact message sent with ID: %s\n", id)
}
