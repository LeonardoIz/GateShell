
# GateShell

*A modern reverse proxy for SSH.* 

## How It Works

GateShell acts as an intermediary between clients and backend SSH servers. 
When a client connects to GateShell, it forwards the connection to the appropriate backend server based on the configuration. 
This allows for centralized management of multiple SSH servers.

## Configuration

The configuration is managed through a JSON file. The default configuration file is `config.json`, but you can specify a different path using a command-line flag or environment variable.

Example:

```json
{
  "server": {
    "port": 22,
    "host_key": "ssh_host_key",
    "default_endpoint": "default"
  },
  "endpoints": {
    "default": {
      "target": "backend_server:22",
      "auth": {
        "user": "username",
        "methods": ["password"]
      }
    }
  }
}
```

- `server`: Configuration for the GateShell server.
  - `port`: Port on which GateShell listens for incoming SSH connections.
  - `host_key`: Path to the SSH host key file.
  - `default_endpoint`: Default endpoint to use if no specific endpoint is matched.
- `endpoints`: Configuration for backend endpoints.
  - `name` The username with which the endpoint will be accessed, you can add as many endpoints as you want.
    - `target`: The address of the backend SSH server.
    - `auth`: Authentication configuration for the backend server.
      - `user`: Username for authentication.
      - `methods`: List of authentication methods, for now only password.

## How to Run

1. Clone and build the repository:

    ```sh
    git clone https://github.com/LeonardoIz/GateShell.git
    cd GateShell
    go build -o gateshell cmd/gateshell
    ```

2. Create a configuration file `config.json` in the root directory. Use the example configuration file structure provided above.

3. Run the GateShell server:

    ```sh
    ./gateshell -config config.json
    ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
