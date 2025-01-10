# Go Chat Application

This project is a simple real-time chat application built using Go, WebAssembly, and WebSockets. It demonstrates a basic architecture for handling concurrent connections and message broadcasting.

## Project Structure

The project is organized as follows:

-   **`.air.toml`**: Configuration file for the `air` live-reloading tool.
-   **`scripts/build_wasm.sh`**: Shell script to build the WebAssembly (WASM) frontend.
-   **`main.go`**: The main server-side Go application.
-   **`frontend/main.go`**: The frontend Go application compiled to WASM.
-   **`.vscode/settings.json`**: VS Code settings for the project.
-   **`.gitignore`**: Specifies intentionally untracked files that Git should ignore.
-   **`frontend/index.html`**: The main HTML file for the frontend.
-   **`shared/types.go`**: Defines shared data structures used by both the server and the frontend.
-   **`dist/`**: Directory where the static frontend files are served from.
-   **`tmp/`**: Temporary directory used during the build process.

### Key Components

-   **Server (`main.go`)**:
    -   Sets up a WebSocket server using `gorilla/websocket`.
    -   Uses `github.com/anthdm/hollywood` for actor-based concurrency.
    -   Handles client connections and message routing.
    -   Implements a `RoomActor` to broadcast messages to all connected clients.
    -   Serves static files from the `dist` directory.
-   **Frontend (`frontend/main.go`)**:
    -   Written in Go and compiled to WebAssembly.
    -   Uses `github.com/hexops/vecty` for UI rendering.
    -   Connects to the WebSocket server.
    -   Displays incoming messages in a chat window.
    -   Allows users to send messages.
-   **Shared (`shared/types.go`)**:
    -   Defines common types like `WSMessage` and `TextMessage` for communication between the server and frontend.

## How to Run

1.  **Install Dependencies:**
    -   Make sure you have Go installed (version 1.18 or later).
    -   Install `air` for live reloading: `go install github.com/cosmtrek/air@latest`
    -   Install `tailwindcss` for styling: `npm install -D tailwindcss`
        -   Initialize Tailwind CSS: `npx tailwindcss init -p`
2.  **Build the Frontend:**
    -   Run the `scripts/build_wasm.sh` script:
        ```bash
        ./scripts/build_wasm.sh
        ```
    -   This will copy `wasm_exec.js` and build the `main.wasm` file in the `frontend/` directory.
3.  **Build the Tailwind CSS:**
    -   Create a `tailwind.config.js` file in the root directory with the following content:
        ```javascript
        /** @type {import('tailwindcss').Config} */
        module.exports = {
          content: [
            "./frontend/index.html",
            "./frontend/main.go",
          ],
          theme: {
            extend: {},
          },
          plugins: [],
        }
        ```
    -   Create a `frontend/input.css` file with the following content:
        ```css
        @tailwind base;
        @tailwind components;
        @tailwind utilities;
        ```
    -   Run the following command to build the CSS:
        ```bash
        npx tailwindcss -i ./frontend/input.css -o ./dist/output.css --watch
        ```
4.  **Start the Server:**
    -   Run `air` from the project root:
        ```bash
        air
        ```
    -   This will start the Go server and watch for file changes, automatically rebuilding and restarting the server when needed.
5.  **Open in Browser:**
    -   Open your web browser and navigate to `http://localhost:8080`.

## Key Concepts

-   **WebSockets:** Used for bidirectional, real-time communication between the server and the browser.
-   **WebAssembly (WASM):** Allows running Go code in the browser.
-   **Actor Model:** Used for concurrency on the server-side, enabling efficient handling of multiple client connections.
-   **Live Reloading:** The `air` tool provides live reloading during development, making it easier to see changes in real-time.
-   **Tailwind CSS:** Used for styling the frontend.

## Notes

-   The `CheckOrigin` in the WebSocket upgrader is set to allow all origins for simplicity. In a production environment, you should configure this to only allow specific origins.
-   The frontend is very basic and can be extended with more features.
-   The `RoomActor` currently broadcasts all messages to all clients. You could implement more sophisticated routing or user-specific messaging.
-   The `main.go` file uses the `log` package from `github.com/charmbracelet/log` for structured logging.

## Contributing

Feel free to contribute to this project by submitting pull requests or opening issues.

## License

This project is licensed under the MIT License.