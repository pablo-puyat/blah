# blah

a log viewer to ease web development

## Log Monitor Architecture

## Overview

Log Monitor is a developer tool designed to enhance local log monitoring by providing real-time log viewing, searching, and aggregation capabilities. It replaces traditional `tail -f` with a more developer-friendly approach while maintaining simplicity and avoiding external dependencies.

## Core Features

- Real-time log file monitoring
- Multiple log source support (Laravel.log, JavaScript errors)
- In-memory buffering for performance
- SQLite-based storage for search capabilities
- WebSocket-based real-time updates
- Single binary distribution

## Architecture

### Components

#### 1. File Watcher Module

**Responsibilities:**

- Monitor specified files for changes
- Read new content
- Distribute entries to processing pipeline

#### 2. Display Module 

**Responsibilities:**

- Listen to log channel
- Output entries to terminal
- Provide TUI to navigate entries
    - paging for individual entries
    - collapse and entries

#### 3. Storage Module

SQLite-based storage system for log persistence and searching.

**Responsibilities:**

- Listen to log channel
- Persist logs to SQLite database
- Maintain search indexes
- Handle historical queries
- Implement log retention policies

#### 4. Summarizer Module

Processes and aggregates log data for analysis.

**Responsibilities:**

- Calculate real-time statistics
- Generate time-based summaries
- Track error patterns
- Provide monitoring metrics

#### 5. WebSocket Server

Handles real-time communication with clients.

**Responsibilities:**

- Manage viewer connections
- Accept frontend log submissions
- Broadcast log updates
- Handle client queries

### Data Flow

```
Log Sources
├── File Watcher (Laravel.log)
└── Frontend Logger (JS Errors)
       │
       ▼
    Channel
       │
       ├─────► Display
       │
       ├─────► SQLite Writer
       │
       ├─────► Summarizer
       │
       └─────► WebSocket Broadcaster
```

## Implementation Details

### Channel Design

- Buffered channels for handling spikes
- Configurable buffer sizes
- Drop policy for overflow scenarios
- Error channel for system issues
- Context preservation for log correlation

### Configuration

```go
type Config struct {
    WatchPaths    []string // Paths to watch
    BufferSize    int      // Channel buffer size
    WebSocketPort int      // WS server port
    DBPath        string   // SQLite database path
    RetentionDays int      // Log retention period
}
```

### Error Handling

The system must handle:

- File system errors (missing files, permissions)
- Channel overflow scenarios
- Database errors
- Client connection issues
- Message parsing failures

## Development Phases

### Phase 1: Core Functionality

- File watching implementation
- Basic stdout output

### Phase 2: Storage Layer

- SQLite integration

### Phase 3: Real-time Updates

- WebSocket server implementation
- Basic client broadcasting
- Connection management

### Phase 4: Frontend Client

- Log viewer interface
- Search functionality
- Real-time updates

### Phase 5: Frontend Integration

- JavaScript error capturing
- Frontend logger implementation

## Memory Management

### Buffer Strategy

- Fixed-size channel buffers
- Rolling window of recent logs
- Query-based historical access
- Automatic cleanup of old entries

### Performance Considerations

- Efficient file reading
- Batched database writes
- Optimized WebSocket broadcasting
- Resource-conscious client updates

## Usage Examples

### Basic Usage

```bash
blah /path/to/laravel.log
```

### Multiple Sources

```bash
blah /path/to/laravel.log /path/to/access.log
```

### With Frontend Logging

```javascript
// Frontend integration
blah.connect({
    port: 8080,
    captureErrors: true
});
```

## Development Guidelines / Goal

1. produce a single binary
2. Maintain single responsibility per module
3. Implement clear interfaces
4. Ensure easy deployment
5. Prioritize user experience

## Future Considerations

- Log format plugins
- Enhanced aggregation features
- Custom alerting rules
