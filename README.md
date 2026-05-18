
# goma

`goma` is a high-performance, pipeline-oriented CLI utility written in Go for sending emails via `msmtp`. 
It specializes in handling dynamic content from `stdin` and managing complex MIME multipart messages for attachments without the overhead of heavy dependencies.

---

## Prerequisites

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

### Setting up a Gmail App Password

Google does not allow applications to authenticate using your standard Google account password. Instead, you must generate a 16-character **App Password**:

1. Go to your [Google Account Console](https://myaccount.google.com/).
2. On the left navigation panel, select **Security**.
3. Under *How you sign in to Google*, ensure **2-Step Verification** is turned **ON** (this is mandatory to use App Passwords).
4. Click on **2-Step Verification**, scroll all the way to the bottom of the page, and select **App passwords**.
5. Enter a custom name for your app (e.g., `goma` or `msmtp`) and click **Create**.
6. Copy the generated **16-character code** shown in the yellow box.
7. Paste this code into your `~/.msmtprc` file in place of `xxxx xxxx xxxx xxxx` (**do not include spaces**; `msmtp` expects a solid 16-character string).

**Security Tip:** Secure your `~/.msmtprc` file immediately after adding your password so other local users cannot read it:
```bash
chmod 600 ~/.msmtprc

```


---

## Installation

### 1. Clone and Build

```bash
git clone [https://github.com/mmngadi/goma.git](https://github.com/mmngadi/goma.git)
cd goma
go build -o goma main.go

```

### 2. Install Globally

Move the compiled binary to your local execution path to make it globally accessible:

```bash
sudo mv goma /usr/local/bin/goma

```

---

## Usage

`goma` reads the email body from `stdin`. 

This allows it to sit at the end of a processing chain (e.g., after [gsub](https://github.com/mmngadi/gsub-pub) or `cat`).

### Basic Text Alert

```bash
echo "Backup successful" | goma -to "admin@example.com" -subject "System Alert"

```

### Sending HTML Alerts (with gsub)

```bash
cat alert.html.tmpl | gsub -e | goma \
  -to "noc@example.com" \
  -subject "Node Status: CRITICAL" \
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
