package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Custom type to handle both comma-separated and multiple flags
type attachList []string

func (a *attachList) String() string {
	return strings.Join(*a, ", ")
}

func (a *attachList) Set(value string) error {
	// Split by comma to support -attach "file1.txt,file2.log"
	paths := strings.Split(value, ",")
	for _, p := range paths {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			*a = append(*a, trimmed)
		}
	}
	return nil
}

func main() {
	var toAddr string
	var subject string
	var contentType string
	var account string
	var attachments attachList

	flag.StringVar(&toAddr, "to", "", "Recipient email address (Required)")
	flag.StringVar(&subject, "subject", "Notification", "Email subject line")
	flag.StringVar(&contentType, "type", "text/plain", "Content type: text/plain or text/html")
	flag.StringVar(&account, "account", "gmail", "msmtp account profile to use")
	flag.Var(&attachments, "attach", "Path(s) to file (comma-separated or multiple flags)")

	flag.Parse()

	if toAddr == "" {
		fmt.Fprintln(os.Stderr, "Error: The -to flag is mandatory.")
		os.Exit(1)
	}

	body, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading body: %v\n", err)
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// Headers
	fmt.Fprintf(buf, "To: %s\n", toAddr)
	fmt.Fprintf(buf, "Subject: %s\n", subject)
	fmt.Fprintf(buf, "MIME-Version: 1.0\n")
	fmt.Fprintf(buf, "Content-Type: multipart/mixed; boundary=%s\n\n", writer.Boundary())

	// Body Part
	bodyHeader := make(textproto.MIMEHeader)
	bodyHeader.Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	bodyPart, _ := writer.CreatePart(bodyHeader)
	bodyPart.Write(body)

	// Attachment Parts
	for _, path := range attachments {
		if err := addFile(writer, path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Skipping %s: %v\n", path, err)
		}
	}

	writer.Close()

	// Execute msmtp
	cmd := exec.Command("msmtp", "-a", account, toAddr)
	cmd.Stdin = buf
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "msmtp error: %v\nOutput: %s\n", err, string(output))
		os.Exit(1)
	}

	fmt.Printf("✓ Email sent to %s with %d attachment(s).\n", toAddr, len(attachments))
}

func addFile(writer *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath)))
	header.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	encoder := base64.NewEncoder(base64.StdEncoding, part)
	defer encoder.Close()
	_, err = io.Copy(encoder, file)
	return err
}
