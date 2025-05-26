# 📡 Pager Notification System

Pager is a **basic but scalable notification system architecture** that demonstrates core patterns for handling notifications efficiently. While not production-grade out of the box, it provides the foundational components for:

- Batch processing notifications
- Multi-channel delivery (Kafka, Email, SMS)
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
```mermaid
graph TD
    A[API Server] --> B[Notification Session]
    B --> C[Fetch Audience]
    C --> D[Batch Processor\n(5 audience/batch)]
    D --> E[Create Batch]
    E --> F[Push Batch to Kafka]
    F --> G[Kafka Consumer]
    G --> H[Process Batch]
    H --> I[Fetch Templates]
    I --> J[Personalize Content]
    J --> K[Notification Service\nProvider]
    
    style A fill:#f9f,stroke:#333
    style B fill:#bbf,stroke:#333
    style C fill:#6f9,stroke:#333
    style D fill:#f96,stroke:#333
    style E fill:#9cf,stroke:#333
    style F fill:#ff9,stroke:#333
    style G fill:#c9f,stroke:#333
    style H fill:#f66,stroke:#333
    
    classDef component fill:#fff,stroke:#333,stroke-width:2px;
    class A,B,C,D,E,F,G,H component;
```

The complete batch processing flow:
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
- **Atomic Processing**: Batch succeeds or fails as a unit
- **Kafka Integration**:
  - Producer pushes complete batches
  - Consumer processes one batch at a time
- **Template Personalization**: Done at consumer level
- **Pluggable Providers**: Support for Email/SMS/etc.

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

# 3. Install & run
go mod download
source .env && go run main.go apis
```

## 🧪 Testing the System
```bash
# Run unit tests
go test -v ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out


## 🚀 Deployment Options

---
✨ **Happy Notifying!** ✨
