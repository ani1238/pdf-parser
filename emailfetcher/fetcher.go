package emailfetcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func FetchPdfsFromEmailForSubject(userEmail, subject string) error {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return err
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	messages, err := listMessagesWithSubject(srv, userEmail, subject)
	if err != nil {
		return err
	}
	fmt.Println("Messages:", messages)

	for _, message := range messages {
		fmt.Println("Message ID:", message.Id)

		// Optionally, retrieve and decode the full content of each email
		fullMessage, err := GetMessage(srv, userEmail, message.Id)
		if err != nil {
			log.Printf("Error retrieving message: %v", err)
			continue
		}
		pdfAttachments := make([]*gmail.MessagePart, 0)
		for _, part := range fullMessage.Payload.Parts {
			if part.MimeType == "application/pdf" {
				pdfAttachments = append(pdfAttachments, part)

			}
			// fmt.Println("mime type ", part.MimeType)

		}
		if len(pdfAttachments) > 0 {
			for _, pdfAttachment := range pdfAttachments {
				attachment, err := GetAttachment(srv, userEmail, message.Id, pdfAttachment.Body.AttachmentId)
				if err != nil {
					log.Printf("Error retrieving attachment: %v", err)
					continue
				}
				decodedData, _ := DecodeBase64(attachment.Data)
				err = saveAttachment(pdfAttachment.Filename, "outpdfs", decodedData)
				if err != nil {
					log.Printf("Error saving attachment: %v", err)
					continue
				}
			}

		}

	}
	return nil
}

func saveAttachment(filename, destDir string, data []byte) error {
	filePath := fmt.Sprintf("%s/%s", destDir, filename)
	err := writeFile(filePath, data)
	if err != nil {
		return err
	}
	log.Printf("Attachment saved: %s\n", filePath)
	return nil
}

func writeFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func DecodeBase64(body string) ([]byte, error) {
	data, err := base64.URLEncoding.DecodeString(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetMessage(srv *gmail.Service, user, msgID string) (*gmail.Message, error) {

	message, err := srv.Users.Messages.Get(user, msgID).Do()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func GetAttachment(srv *gmail.Service, user, msgID string, attachmentId string) (*gmail.MessagePartBody, error) {

	attachment, err := srv.Users.Messages.Attachments.Get(user, msgID, attachmentId).Do()
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func listMessagesWithSubject(srv *gmail.Service, user, subject string) ([]*gmail.Message, error) {
	query := fmt.Sprintf("subject:%s", subject)

	response, err := srv.Users.Messages.List(user).Q(query).Do()
	if err != nil {
		return nil, err
	}

	return response.Messages, nil
}
