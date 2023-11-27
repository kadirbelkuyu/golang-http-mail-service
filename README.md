# Go Mail Service

This is a simple mail service written in Go.

## Requirements

- Go 1.16 or later
- An SMTP server

## Installation

Clone the repository:

```bash
git clone https://github.com/kadirbelkuyu/golang-http-mail-service.git
```
Navigate to the project directory:

```bash
cd gomailservice
```

Install the dependencies:
```bash
go mod download
```

```bash
docker compose up -d 
```


Create a `.env` file in the root directory and fill it with your SMTP server details:

```bash
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SENDER_EMAIL=your-email@example.com
SENDER_PASSWORD=your-email-password
```

### Running the Application
Run the application:

```bash
go run cmd/server/main.go
```

## Usage
To send an email, make a POST request to `/sendmail` with the following JSON body:

```json
{
    "to": "recipient@example.com",
    "subject": "Hello",
    "body": "Hello, world!"
}
``````

License
This project is licensed under the MIT License.
---
