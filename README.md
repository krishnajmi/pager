# 📡 Pager Notification System

Pager is a **basic but scalable notification system architecture** that demonstrates core patterns for handling notifications efficiently. While not production-grade out of the box, it provides the foundational components for:

- Batch processing notifications
- Basic retry mechanisms
- Monitoring capabilities

The architecture is intentionally kept simple to serve as:
1. A learning resource for Go distributed systems
2. A starting point for more complex implementations
3. A reference for notification system design patterns

## 🌟 Key Features

### 🚀 Core Capabilities
| Feature | Description |
|---------|-------------|
| **Batch Processing** | Efficiently processes notifications in configurable batches (default: 5 per batch) |
| **Template Processing** | Powerful template engine with dynamic personalization, variable substitution, and conditional content |
| **Authentication & Authorization** | Role-based access control (RBAC) with granular permissions for notification management |
| **Security** | End-to-end encryption, secure credential storage, and audit logging |

### 🧩 Pager Architecture

The complete notification processing flow:
1. **API Server** receives the notification request
2. Creates a **Notification Session** to track the process
3. **Fetches Audience** data (from request in this implementation)
4. **Batch Processor** groups audiences in batches of 5
5. **Creates Batch** messages containing all audience data
6. **Pushes complete Batch** to Kafka queue
7. **Kafka Consumer** pulls batches from the queue
8. **Processes each Batch** sequentially
9. For each batch item: **Fetches Template** based on type
10. **Personalizes Content** for each recipient
11. Sends to final **Notification Service Provider**

Key Batch Processing Characteristics:
- **Batch-First Approach**: Entire batches processed together
- **Configurable Batch Size**: Default 5 audiences per batch
- **Kafka Integration**:
  - Producer pushes complete batches
  - Consumer processes one batch at a time
- **Template Personalization**: Done at consumer level
- **Pluggable Providers**: Can add Support for multiple channels like Email/SMS/etc.

## 🛠️ Getting Started

### 📋 Prerequisites
- Go 1.18+
- PostgreSQL 12+ (for user data)
- Kafka 2.8+ (for message queue)
- Redis (for caching)

### 🚀 Quick Installation
```bash
# 1. Clone repository
git clone https://github.com/kp/pager.git && cd pager

# 2. Setup environment
cp .env.example .env
vim .env  # Configure your settings

# 3. Start kafka dependencies
docker-compose -f kafka/kafka-compose.yml up -d

# 4. Install & run
go mod download
source .env && go run main.go apis

# migrate db separately
./pager migrate

# 4.2 Register User 
go build
./pager register -u username -p password
```


## 🧪 Testing the System
```bash
# Run unit tests
go test -v ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

## API Documentation

### Authentication
1. First login to get an X-Auth-Token:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'
```

2. Use the returned X-Auth-Token for subsequent requests:



## 🚀 Deployment Options
Orchestration and Docker deployment is coming soon

---
✨ **Happy Notifying!** ✨
