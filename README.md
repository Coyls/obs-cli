# obs-cli

A command-line tool to manage your Obsidian vault.

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Building

1. Clone the repository:

```bash
git clone https://github.com/coyls/obs-cli.git
cd obs-cli
```

2. Build the project:

```bash
chmod +x build.sh
./build.sh
```

Binaries will be generated in the `bin/` directory:

- `obs-cli-linux-amd64` for Linux
- `obs-cli-windows-amd64.exe` for Windows
- `obs-cli-darwin-amd64` for macOS

3. Copy the binary to your PATH:

```bash
# For Linux/macOS
sudo cp bin/obs-cli-linux-amd64 /usr/local/bin/obs-cli

# For Windows
# Copy obs-cli-windows-amd64.exe to a directory in your PATH
```

## Configuration

1. Initialize the configuration:

```bash
obs-cli init
```

2. The configuration file will be created at `~/.config/obs-cli/config.yaml`

3. Modify the configuration according to your needs:

```yaml
vault_path: "/path/to/your/vault"
default_mv_path: "Archives"
default_cp_path: "Assets"
default_editor: "code" # Default editor for editing files (e.g., "code" for VS Code)
```

## Usage

### Available Commands

- `obs-cli init` : Initialize or update the configuration
- `obs-cli mv [file]` : Move a file to the vault
- `obs-cli cp [file]` : Copy a file to the vault
- `obs-cli push` : Push changes to GitHub
- `obs-cli pull` : Pull changes from GitHub
- `obs-cli callouts` : Edit Obsidian callouts configuration

### Examples

```bash
# Move a file to the vault
obs-cli mv ~/Downloads/image.png -d Assets/new

# Copy a file to the vault
obs-cli cp ~/Pictures/photo.jpg -d Assets/new

# Use the default directory
obs-cli mv ~/Documents/document.pdf

# Edit Obsidian callouts
obs-cli callouts
```

## License

MIT
