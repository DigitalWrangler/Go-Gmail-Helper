package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Load OAuth 2.0 credentials
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Read token from file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Get a new token via OAuth2 authentication
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Open this link and authenticate: \n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the code: ")
	fmt.Scan(&authCode)

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Could not get token: %v", err)
	}
	return tok
}

// Save token to file
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving token file: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Could not create token file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Fetch ALL unread emails using pagination
func getAllUnreadEmails(srv *gmail.Service, user string) ([]*gmail.Message, error) {
	var messages []*gmail.Message
	pageToken := ""

	for {
		req := srv.Users.Messages.List(user).Q("is:unread").MaxResults(500) // Fetch up to 500 per request
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		resp, err := req.Do()
		if err != nil {
			return nil, err
		}

		messages = append(messages, resp.Messages...)

		// If no more pages, break loop
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return messages, nil
}

func main() {
	ctx := context.Background()

	// Read credentials.json
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Could not read credentials.json: %v", err)
	}

	// Create OAuth2 configuration with Modify scope
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Could not create OAuth2 config: %v", err)
	}

	client := getClient(config)

	// Create Gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Could not create Gmail service: %v", err)
	}

	user := "me"

	// Get ALL unread emails
	unreadMessages, err := getAllUnreadEmails(srv, user)
	if err != nil {
		log.Fatalf("Could not fetch messages: %v", err)
	}

	totalUnread := len(unreadMessages)
	fmt.Printf("📩 Total unread emails: %d\n", totalUnread)

	if totalUnread == 0 {
		fmt.Println("✅ No unread messages found.")
		return
	}

	markedAsRead := 0 // Counter for emails marked as read

	// Process each unread email
	for _, m := range unreadMessages {
		msg, err := srv.Users.Messages.Get(user, m.Id).Format("metadata").Do()
		if err != nil {
			log.Printf("Could not retrieve message %s: %v", m.Id, err)
			continue
		}

		var subject string
		var dateReceived time.Time

		// Extract headers (Subject, Date)
		for _, header := range msg.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
			}
			if header.Name == "Date" {
				dateReceived, _ = time.Parse(time.RFC1123Z, header.Value)
			}
		}

		// Print email details
		fmt.Printf("📧 %s - Received: %s\n", subject, dateReceived.Format("2006-01-02 15:04:05"))

		// Mark as read
		_, err = srv.Users.Messages.Modify(user, m.Id, &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{"UNREAD"},
		}).Do()
		if err != nil {
			log.Printf("❌ Failed to mark message %s as read: %v", m.Id, err)
		} else {
			fmt.Println("✅ Marked as read!")
			markedAsRead++
		}
	}

	// Fetch unread count after marking
	remainingUnread, err := getAllUnreadEmails(srv, user)
	if err != nil {
		log.Fatalf("Could not fetch updated unread count: %v", err)
	}

	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("🔹 Emails marked as read: %d\n", markedAsRead)
	fmt.Printf("📩 Remaining unread emails: %d\n", len(remainingUnread))
}
