# Glyra API Reference

Glyra establishes a bridge between the Go backend and the Javascript frontend using a native webview implementation.

## Calling Go from JavaScript

In your Go application, you can bind functions to the `window` object in JavaScript using `w.Bind()`.

### Backend (Go)
```go
w.Bind("FetchUserData", func(id int) string {
    // Connect to database, fetch user
    return fmt.Sprintf("User %d fetched from Go!", id)
})
```

### Frontend (JavaScript/TypeScript)
```javascript
// The function is automatically exposed globally on the window object
// It always returns a Promise!
const data = await window.FetchUserData(123);
console.log(data); // "User 123 fetched from Go!"
```

## Security
Because the frontend runs inside a native webview, it has full access to the `window` bindings without relying on HTTP or WebSockets, avoiding CORS issues and network overhead.
