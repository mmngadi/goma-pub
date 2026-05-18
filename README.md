
`goma` is a high-performance, pipeline-oriented CLI utility written in Go for sending emails via `msmtp`. 
It specializes in handling dynamic content from `stdin` and managing complex MIME multipart messages for attachments without the overhead of heavy dependencies.

---

## 🛠 Prerequisites

Before using `goma`, ensure your environment meets the following requirements:

1. **Go Toolchain**: [Installed and configured](https://go.dev/doc/install).
2. **msmtp**: A lightweight SMTP client.
   - **Fedora**: `sudo dnf install msmtp`
   - **Ubuntu/Debian**: `sudo apt install msmtp`
3. **msmtp Configuration**: Ensure `~/.msmtprc` is configured with at least one account (e.g., `gmail`).
   ```bash
   # Example ~/.msmtprc
   defaults
   auth           on
   tls            on
   tls_trust_file /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem
   logfile        ~/.msmtp.log
   
   account        gmail
   host           smtp.gmail.com
   port           587
   from           yourname@gmail.com
   user           yourname@gmail.com
   password       xxxx xxxx xxxx xxxx
   ```
---

## Installation

### 1. Clone and Build

```bash
git clone https://github.com/mmngadi/goma.git
cd goma
go build -o goma main.go

```

### 2. Install Globally

Create a symbolic link to `/usr/local/bin` to call `goma` from anywhere:

```bash
sudo ln -s $(pwd)/goma /usr/local/bin/goma

```

---

## Usage

`goma` reads the email body from `stdin`. This allows it to sit at the end of a processing chain (e.g., after `gsub` or `cat`).

### Basic Text Alert

```bash
echo "Backup successful" | goma -to "admin@example.com" -subject "System Alert"

```

### Sending HTML Alerts (with gsub)

```bash
cat alert.html.tmpl | gsub -e | goma \\
  -to "noc@example.com" \\
  -subject "Node Status: CRITICAL" \\
  -type "text/html"

```

### Advanced Attachments

`goma` features a hybrid attachment system. You can use multiple flags or comma-separated strings.

```bash
# Using explicit multiple flags
goma -to "user@me.com" -attach "logs.txt" -attach "report.pdf" < body.txt

# Using comma-separated lists
goma -to "user@me.com" -attach "diag.log,config.json" < body.txt

```

---

## Configuration Flags

| Flag | Argument | Description | Default |
| --- | --- | --- | --- |
| `-to` | `string` | **Required.** Recipient email address. | N/A |
| `-subject` | `string` | The subject line of the email. | `Notification` |
| `-type` | `string` | Content-Type header (`text/plain`, `text/html`). | `text/plain` |
| `-account` | `string` | The `msmtp` account to use from your config. | `gmail` |
| `-attach` | `string` | File paths (repeatable or comma-separated). | N/A |

---

## License

Distributed under the MIT License.
