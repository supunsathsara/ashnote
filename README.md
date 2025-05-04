# AshNote 🔥

AshNote is a secure "Burn After Reading" message application that allows users to create encrypted messages that self-destruct after being read once.

## Features

- **End-to-End Encryption**: Messages are encrypted with AES using a password provided by the sender
- **One-Time Access**: Messages are permanently deleted after being read once
- **User-Friendly Interface**: Clean, responsive design for easy use
- **Secure-by-Design**: No plaintext storage of sensitive data

## Project Structure

```
ashnote/
├── cmd/
│   └── main.go                # Application entry point
├── internal/
│   ├── crypto/
│   │   └── crypto.go          # Encryption and decryption utilities
│   ├── db/
│   │   └── db.go              # Database operations
│   └── handler/
│       └── handler.go         # HTTP request handlers
├── web/
│   └── templates/
│       ├── index.html         # Homepage template
│       └── message.html       # Message viewing template
├── Dockerfile                 # Docker configuration
├── README.md                  # This file
├── go.mod                     # Go module dependencies
└── go.sum                     # Go module checksums
```

## Prerequisites

- Go 1.23 or higher
- SQLite (included as a Go dependency)

## Running Locally

1. Clone the repository:

   ```bash
   git clone https://github.com/supunsathsara/ashnote.git
   cd ashnote
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run cmd/main.go
   ```

4. Open your browser and navigate to [http://localhost:3000](http://localhost:3000)

## Docker Deployment

1. Build the Docker image:

   ```bash
   docker build -t ashnote .
   ```

2. Run the container:

   ```bash
   docker run -p 3000:3000 -v ashnote-data:/app/data ashnote
   ```

3. Access the application at [http://localhost:3000](http://localhost:3000)

## How It Works

1. The sender creates a message and provides a password
2. The message is encrypted using AES encryption with the password
3. The encrypted message is stored in a SQLite database with a unique ID
4. A unique URL is generated for accessing the message
5. When the recipient opens the URL, they are prompted for the password
6. If the password is correct, the message is decrypted and displayed
7. After successful decryption, the message is permanently deleted from the database

## Security Considerations

- Messages are encrypted at rest using AES encryption
- Passwords are never stored - they are only used for encryption/decryption
- Messages are automatically deleted after being read once
- SQLite database file should be properly secured in production environments

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
