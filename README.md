# URL Shortener with Go

[![License](https://img.shields.io/github/license/ubayed-bin-sufian/url-shortener-with-go)](LICENSE)

A lightning-fast, production-ready URL shortener service built with Go and Redis. Transform long URLs into memorable, short links with custom aliases support.

## Features

- **Instant URL Shortening**: Create short URLs in milliseconds
- **Custom Aliases**: Choose your own custom short URL aliases
- **Rate Limiting**: Built-in protection against abuse
- **Redis Backend**: High-performance data storage

## Architecture

The service uses:
- Go Fiber for the HTTP server
- Redis for storing URL mappings and rate limiting
- Docker for containerization
- Two Redis databases:
  - DB 0: URL mappings
  - DB 1: Rate limiting data

## Quick Start

### Prerequisites

- Go 1.24.4 or higher
- Redis 6.x or higher
- Docker and Docker Compose

### Installation

1. Clone the repository
```bash
git clone https://github.com/ubayed-bin-sufian/url-shortener-with-go.git
cd url-shortener-with-go
```

2. Install dependencies
```bash
go mod download
```

### Environment Variables

Create a `.env` file in the `api` directory with these variables:
```env
DB_ADDR="db:6379"    # Redis address
DB_PASS=""           # Redis password
APP_PORT=":3000"     # Application port
DOMAIN="localhost:3000"  # Your domain
API_QUOTA=10         # Rate limit per IP
```

### Docker Setup

```bash
docker-compose up -d
```

## API Usage

### Create a short URL

```bash
curl -i -X POST -H "Content-Type: application/json" -d '{
    "url": "https://github.com",
    "expiry": 24
}' http://localhost:3000/api/v1
```

Example Response:
```json
{
    "url": "https://github.com",
    "short": "localhost:3000/abc123",
    "expiry": 24,
    "rate-limit": 9,
    "rate-limit-reset": 30
}
```

Access the shortened URL:

```bash
curl -i http://localhost:3000/your-short-code
```

### Use custom short URL

```bash
curl -i -X POST -H "Content-Type: application/json" -d '{
    "url": "https://github.com",
    "short": "gh3",
    "expiry": 24
}' http://localhost:3000/api/v1
```

To resolve/use a shortened URL:

```bash
curl -i http://localhost:3000/gh3
```

Each IP address gets 10 requests (configured in .env as `API_QUOTA`). The rate limit resets after 30 minutes.

You can also use your web browser to test the shortened URLs by visiting `http://localhost:3000/your-short-code`.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
