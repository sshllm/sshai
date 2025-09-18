# SSHAI - SSH AI Assistant

English | [简体中文](./README.md)

An intelligent AI assistant program that provides AI model services through SSH connections, allowing you to use AI assistants in any SSH-supported environment.

## ✨ Key Features

- 🔐 **Secure SSH Connection** - Encrypted AI service access through SSH protocol
- 🔑 **Flexible Authentication** - Support for password authentication, SSH key-based passwordless login, and passwordless mode
- 🗝️ **SSH Keys Support** - Support for multiple SSH public keys for passwordless login, compatible with RSA, Ed25519, and other key types
- 🤖 **Multi-Model Support** - Support for DeepSeek, Hunyuan, and other AI models
- 💭 **Real-time Thinking Display** - Real-time display of thinking processes for models like DeepSeek R1
- 🎨 **Beautiful Interface** - Colorful output, animations, and ASCII art
- ⚙️ **Flexible Configuration** - Support for dynamic configuration file specification (-c parameter) and complete YAML configuration
- 🌐 **Multi-language Support** - Support for Chinese and English interfaces
- 📝 **Custom Prompts** - Configurable AI prompt system
- 🚀 **Startup Welcome Banner** - Beautiful welcome banner displayed on program startup
- 🏗️ **Modular Design** - Clean code architecture, easy to extend

## 🚀 Quick Start

### 1. Download and Build

```bash
# Clone the project
git clone https://github.com/sshllm/sshai.git
cd sshai

# Build the program
make build
# or
go build -o sshai cmd/main.go
```

### 2. Configuration Setup

Edit the `config.yaml` file and set your API key:

```yaml
# API Configuration
api:
  base_url: "https://api.deepseek.com/v1"
  api_key: "your-api-key-here"
  default_model: "deepseek-v3"

# Server Configuration
server:
  port: 2213
  welcome_message: "Welcome to SSHAI!"

# Authentication Configuration (Optional)
auth:
  password: ""  # Empty = no password authentication
  login_prompt: "Please enter password: "
  # SSH public key passwordless login configuration (only effective when password is set)
  authorized_keys:
    - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user@hostname"
    - "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user2@hostname"
  authorized_keys_file: "~/.ssh/authorized_keys"  # Optional: read public keys from file

# Custom Prompt Configuration
prompt:
  system_prompt: "You are a professional AI assistant, please answer questions in English."
  stdin_prompt: "Please analyze the following content and provide relevant help or suggestions:"
  exec_prompt: "Please answer the following question or execute the following task:"
```

### 3. Run the Server

```bash
# Run directly (using default config.yaml)
./sshai

# Specify configuration file
./sshai -c config.yaml
./sshai -c /path/to/your/config.yaml

# Run in background
./sshai > server.log 2>&1 &

# Run with script
./scripts/run.sh
```

#### Command Line Parameters

- `-c <config_file>` - Specify configuration file path
  - If not specified, defaults to `config.yaml` in current directory
  - If configuration file doesn't exist, program will show error message and exit

```bash
# Usage examples
./sshai -c config.yaml              # Use config file in current directory
./sshai -c /etc/sshai/config.yaml   # Use config file with absolute path
./sshai                             # Default to config.yaml
```

### 4. Connect and Use

```bash
# Interactive mode
ssh user@localhost -p 2213

# Direct command execution
ssh user@localhost -p 2213 "Hello, please introduce yourself"

# Pipe input analysis
cat file.txt | ssh user@localhost -p 2213
echo "Analyze this code" | ssh user@localhost -p 2213
```

## 📁 Project Structure

```
sshai/
├── README.md              # Chinese documentation
├── README_EN.md           # English documentation
├── LICENSE                # Open source license
├── config.yaml           # Main configuration file
├── config-en.yaml        # English configuration file
├── go.mod                # Go module dependencies
├── Makefile              # Build script
├── cmd/                  # Program entry
│   └── main.go           # Main program file
├── pkg/                  # Core modules
│   ├── config/           # Configuration management
│   ├── models/           # Data models
│   ├── ai/               # AI assistant functionality
│   ├── ssh/              # SSH server
│   └── utils/            # Utility functions
├── docs/                 # Project documentation
├── scripts/              # Test and run scripts
└── keys/                 # SSH key files
```

## 🔧 Configuration Guide

### API Configuration

Support for multiple API endpoint configurations:

```yaml
api:
  base_url: "https://api.deepseek.com/v1"
  api_key: "your-deepseek-key"
  default_model: "deepseek-v3"
  timeout: 600

# Multiple API configuration
api_endpoints:
  - name: "deepseek"
    base_url: "https://api.deepseek.com/v1"
    api_key: "your-key"
    default_model: "deepseek-v3"
  - name: "local"
    base_url: "http://localhost:11434/v1"
    api_key: "ollama"
    default_model: "gemma2:27b"
```

### Authentication Configuration

#### Password Authentication
```yaml
auth:
  password: "your-secure-password"  # Set access password
  login_prompt: "Please enter password: "
```

#### SSH Public Key Passwordless Login
```yaml
auth:
  password: "your-secure-password"  # Password must be set to enable SSH public key authentication
  login_prompt: "Please enter password: "
  # Method 1: Configure public key list directly
  authorized_keys:
    - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user@hostname"
    - "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user2@hostname"
  # Method 2: Read public keys from file
  authorized_keys_file: "~/.ssh/authorized_keys"
```

**SSH Public Key Usage**:
```bash
# Generate SSH key pair
ssh-keygen -t ed25519 -f ~/.ssh/sshai_key

# Connect using private key (passwordless login)
ssh -i ~/.ssh/sshai_key -p 2213 user@localhost

# View public key content (for configuration)
cat ~/.ssh/sshai_key.pub
```

**Note**: 
- SSH public key authentication is only enabled when a password is set, providing additional security
- Supports multiple public keys simultaneously, compatible with RSA, Ed25519, ECDSA, and other key types
- After successful login, the program will automatically display built-in welcome information, no need to configure in the config file

### Prompt Configuration

```yaml
prompt:
  system_prompt: "You are a professional AI assistant..."
  stdin_prompt: "Please analyze the following content:"
  exec_prompt: "Please answer the following question:"
```

## 🧪 Testing

The project includes comprehensive test scripts:

```bash
# Basic functionality test
./scripts/test.sh

# SSH execution functionality test
./scripts/test_ssh_exec_final.sh

# Standard input functionality test
./scripts/test_stdin_feature.sh

# Authentication functionality test
./scripts/test_auth.sh

# DeepSeek R1 thinking mode test
./scripts/test_deepseek_r1.sh

# SSH Keys passwordless login functionality test
./scripts/test_ssh_keys.sh
```

## 📚 Documentation

- [Configuration Guide](docs/CONFIG_GUIDE.md) - Detailed configuration instructions
- [Usage Guide](docs/USAGE.md) - Feature introduction and usage methods
- [Architecture Documentation](docs/MODULAR_ARCHITECTURE.md) - Modular architecture design
- [Authentication Configuration](docs/AUTH_CONFIG_EXAMPLE.md) - SSH authentication configuration examples
- [SSH Keys Guide](docs/SSH_KEYS_GUIDE.md) - SSH public key passwordless login configuration guide

## 🤝 Contributing

Issues and Pull Requests are welcome!

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Thanks to all developers and users who have contributed to this project!

---

**Note**: This project follows the Apache 2.0 open source license, welcoming both personal and commercial use.