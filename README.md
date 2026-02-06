# sourcecraft-sdk

A Go client library for the [SourceCraft](https://sourcecraft.dev) Public API.

## About SourceCraft

SourceCraft is a comprehensive platform for software development lifecycle management. It provides:

- **Code Repository Management** - Git-based version control with branching, pull requests, and code review
- **CI/CD Pipeline** - Automated testing and deployment workflows
- **Error Tracking** - Built-in issue tracking and monitoring
- **Code Assistant** - AI-powered code editing, review, and vulnerability analysis
- **Security** - Vulnerability scanning, secret detection, and dependency analysis
- **Package Management** - Support for Maven, NPM, NuGet, PyPI, Docker registries
- **Access Control** - Organization-based permissions and team collaboration

Learn more at [https://sourcecraft.dev/portal/docs/en/](https://sourcecraft.dev/portal/docs/en/)

## About the API

The SourceCraft REST API enables automation and integration with SourceCraft resources including repositories, pull requests, issues, organizations, and more.

### Authentication

The API uses **Personal Access Tokens (PAT)** for authentication. Tokens are passed via the `Authorization: Bearer` header.

For CI/CD workflows, use the `SOURCECRAFT_TOKEN` predefined environment variable.

### API Documentation

- **API Getting Started Guide**: [https://sourcecraft.dev/portal/docs/en/sourcecraft/operations/api-start](https://sourcecraft.dev/portal/docs/en/sourcecraft/operations/api-start)
- **OpenAPI Documentation**: [https://api.sourcecraft.tech/docs/index.html](https://api.sourcecraft.tech/docs/index.html)

## Installation

```bash
go get github.com/aalexzy/sourcecraft-sdk
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    sourcecraft "github.com/aalexzy/sourcecraft-sdk"
)

func main() {
    // Create a new client
    client, err := sourcecraft.NewClient(
        "https://api.sourcecraft.tech",
        sourcecraft.SetToken("your-personal-access-token"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    
    // Use the client to interact with the API
    // Example: List repositories, manage pull requests, etc.
}
```

### Client Options

#### With Personal Access Token

```go
client, err := sourcecraft.NewClient(
    "https://api.sourcecraft.tech",
    sourcecraft.SetToken("your-pat-token"),
)
```

#### With Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: time.Second * 30,
}

client, err := sourcecraft.NewClient(
    "https://api.sourcecraft.tech",
    sourcecraft.SetToken("your-pat-token"),
    sourcecraft.SetHTTPClient(httpClient),
)
```

#### Skip TLS Verification (Development Only)

```go
client, err := sourcecraft.NewClient(
    "https://api.sourcecraft.tech",
    sourcecraft.SetToken("your-pat-token"),
    sourcecraft.WithHTTPClient(true), // insecure=true
)
```

## Webhook Support

The SDK provides comprehensive webhook parsing and verification capabilities, allowing you to easily handle SourceCraft webhook events in your applications.

### Setting Up a Webhook Handler

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    sourcecraft "github.com/aalexzy/sourcecraft-sdk"
)

func main() {
    // Create a webhook parser with secret for HMAC verification
    hook, err := sourcecraft.New(
        sourcecraft.Options.Secret("your-webhook-secret"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Set up HTTP handler
    http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
        payload, err := hook.Parse(r, sourcecraft.PushEvent, sourcecraft.PullRequestAggregate)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Handle the event
        switch p := payload.(type) {
        case sourcecraft.PushEventPayload:
            handlePushEvent(p)
        case sourcecraft.PullRequestEventPayload:
            handlePullRequestEvent(p)
        }

        w.WriteHeader(http.StatusOK)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Handling Push Events

```go
func handlePushEvent(payload sourcecraft.PushEventPayload) {
    fmt.Printf("Push event received for repository: %s\n", payload.Repository.Name)
    fmt.Printf("Pushed by: %s\n", payload.Pusher.Slug)
    fmt.Printf("Branch: %s\n", payload.RefUpdate.RefName)
    fmt.Printf("Default branch updated: %v\n", payload.IsDefaultBranchUpdated)
    
    // Process commits
    for _, commit := range payload.Commits {
        fmt.Printf("  Commit %s by %s: %s\n", 
            commit.Sha[:7], 
            commit.Author.Name, 
            commit.Message)
    }
}
```

### Handling Pull Request Events

The SDK provides a convenient `PullRequestEventPayload` interface that all PR events implement, allowing you to handle all PR events uniformly:

```go
func handlePullRequestEvent(payload sourcecraft.PullRequestEventPayload) {
    pr := payload.GetPullRequest()
    repo := payload.GetRepository()
    eventType := payload.GetEventType()
    
    fmt.Printf("Pull Request event: %s\n", eventType)
    fmt.Printf("Repository: %s\n", repo.Name)
    fmt.Printf("PR #%d: %s\n", pr.Number, pr.Title)
    fmt.Printf("Status: %s\n", pr.Status)
    fmt.Printf("From: %s -> %s\n", pr.SourceBranch, pr.TargetBranch)
    
    // Handle specific event types
    switch eventType {
    case sourcecraft.PullRequestCreateEvent:
        fmt.Println("New pull request created")
    case sourcecraft.PullRequestMergeEvent:
        fmt.Println("Pull request merged")
    case sourcecraft.PullRequestPublishEvent:
        fmt.Println("Pull request published")
    }
}
```

### Handling Specific Pull Request Events

You can also handle specific PR event types individually:

```go
// Parse only specific events
payload, err := hook.Parse(r, 
    sourcecraft.PullRequestCreateEvent,
    sourcecraft.PullRequestMergeEvent,
)
if err != nil {
    log.Printf("Error parsing webhook: %v", err)
    return
}

// Type assert to specific payload type
switch p := payload.(type) {
case sourcecraft.PullRequestCreateEventPayload:
    fmt.Printf("New PR created: %s\n", p.PullRequest.Title)
    fmt.Printf("Created at: %s\n", p.CreatedAt)
    
case sourcecraft.PullRequestMergeEventPayload:
    fmt.Printf("PR merged: %s\n", p.PullRequest.Title)
    fmt.Printf("Merged SHA: %s\n", p.MergeSha)
}
```

### Using Event Aggregates

Event aggregates allow you to subscribe to all events of a certain type without listing each one:

```go
// Listen to all pull request events using aggregate
payload, err := hook.Parse(r, sourcecraft.PullRequestAggregate)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

// All PR events implement PullRequestEventPayload interface
if prPayload, ok := payload.(sourcecraft.PullRequestEventPayload); ok {
    handlePullRequestEvent(prPayload)
}
```

### Supported Events

#### Individual Events

- `PingEvent` - Webhook ping/test event
- `PushEvent` - Git push to repository
- `PullRequestCreateEvent` - Pull request created
- `PullRequestUpdateEvent` - Pull request updated
- `PullRequestPublishEvent` - Pull request published
- `PullRequestRefreshEvent` - Pull request refreshed
- `PullRequestMergeEvent` - Pull request merged
- `PullRequestMergeFailureEvent` - Pull request merge failed
- `PullRequestNewIterationEvent` - New iteration added to PR
- `PullRequestReviewAssignmentEvent` - Reviewer assigned to PR
- `PullRequestReviewDecisionEvent` - Review decision made on PR

#### Event Aggregates

- `PullRequestAggregate` - Matches all `pull_request.*` events

### HMAC Signature Verification

The webhook parser automatically verifies HMAC signatures when a secret is provided:

```go
// Single secret
hook, err := sourcecraft.New(
    sourcecraft.Options.Secret("your-webhook-secret"),
)

// Multiple secrets for secret rotation (comma-separated)
hook, err := sourcecraft.New(
    sourcecraft.Options.Secret("secret1,secret2,secret3"),
)
```

The parser will verify the `X-Src-Signature` header against the request payload using SHA-256 HMAC.

### Without Signature Verification

If you don't need signature verification (not recommended for production):

```go
hook, err := sourcecraft.New()
// No secret option provided
```

## Features

This SDK provides Go bindings for the SourceCraft API, including:

- Organization management
- Repository operations
- Pull request handling
- Webhook parsing and verification with HMAC signature support
- Issue tracking
- Thread-safe client with concurrent access support

## Development

### Running Tests

```bash
go test ./...
```

### Requirements

- Go 1.25.5 or later

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Links

- [SourceCraft Portal](https://sourcecraft.dev)
- [SourceCraft Documentation](https://sourcecraft.dev/portal/docs/en/)
- [API Documentation](https://api.sourcecraft.tech/docs/index.html)
- [API Getting Started](https://sourcecraft.dev/portal/docs/en/sourcecraft/operations/api-start)
