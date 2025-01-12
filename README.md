# Adding New Communications

This guide explains how to add new communication endpoints for both REST API and WebSocket in the codebase.

## Adding a New REST Endpoint

1. Define the Request/Response Types
```go
// In shared/api/types.go or a new type file
type NewFeatureRequest struct {
    // Add request fields
    Field string `json:"field"`
}

type NewFeatureResponse struct {
    // Add response fields
    Result string `json:"result"`
}
```

2. Create a New Route
```go
// In shared/api/routes.go
var NewFeatureRoute = http.NewRoute[NewFeatureRequest, NewFeatureResponse](
    "/api/new-feature",
    http.MethodPost,
)
```

3. Implement the Handler (Server-side)
```go
// In your handler file
func HandleNewFeature(req NewFeatureRequest) (NewFeatureResponse, error) {
    // Implement your handler logic
    return NewFeatureResponse{
        Result: "processed",
    }, nil
}
```

4. Add Client-side Implementation
```go
// In shared/http/client.go or your frontend code
func CallNewFeature(req NewFeatureRequest) (NewFeatureResponse, error) {
    return client.Do(NewFeatureRoute, req)
}
```

## Adding a New WebSocket Message

1. Define the Message Type
```go
// In shared/ws/messages.go
const (
    // Add your new message type
    TypeNewFeature MessageType = "NEW_FEATURE"
)
```

2. Create the Payload Structure
```go
// In shared/ws/messages.go
type NewFeaturePayload struct {
    // Define your payload fields
    Data    string `json:"data"`
    UserID  string `json:"user_id"`
}
```

3. Handle the Message (Server-side)
```go
// In your handler (e.g., internal/actors/client.go)
switch msg.Type {
case ws.TypeNewFeature:
    payload := msg.Payload.(NewFeaturePayload)
    // Handle the new feature message
}
```

4. Send Messages (Client-side)
```go
// In your frontend code
connection.Send(ws.Message{
    Type: ws.TypeNewFeature,
    Payload: NewFeaturePayload{
        Data: "example",
        UserID: "user123",
    },
})
```

## Best Practices

1. **Type Safety**: Always use strongly typed structures for both REST and WebSocket communications.
2. **Documentation**: Add comments explaining the purpose of new message types and payloads.
3. **Validation**: Implement proper validation for all incoming data.
4. **Error Handling**: Define appropriate error responses and handle them gracefully.
5. **Testing**: Add tests for new endpoints and message handlers.

## Common Patterns

- REST endpoints use the `Route[Req, Res]` generic type for type safety
- WebSocket messages follow the `Message` structure with specific payload types
- All WebSocket message types are defined as constants
- Each message type has its own payload structure