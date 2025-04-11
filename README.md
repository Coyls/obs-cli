# obs-cli

A command-line tool to manage my Obsidian vaults.

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
```

## Configuration

1. Set the environment variable for the configuration file:

```bash
# Example
export OBS_CLI_CONFIG="/home/user/my-vaults/.obsclirc.yml"
```

2. Create and modify the configuration file according to your needs:

```yaml
config:
  default_editor: nano # Default editor for editing files (e.g., "code" for VS Code)
  default_vault: MyVault # Name of your default vault
  root: /home/user/my-vaults # Root directory containing your vaults
  vaults:
    MyVault: # Vault name (case-insensitive)
      vault_path: /MyVault # Path to the vault relative to root
      commands:
        cp:
          default_target_path: /assets/new # Default destination for copy command
        mv:
          default_target_path: /assets/new # Default destination for move command
  archive:
    usb_path: /path/to/usb # Path to USB drive for archiving
    extract_path: /path/to/extract # Path where to extract archived files
```

## Usage

### Available Commands

- `obs-cli mv [file]` : Move a file to the vault
- `obs-cli cp [file]` : Copy a file to the vault
- `obs-cli push` : Push changes to GitHub
- `obs-cli pull` : Pull changes from GitHub
- `obs-cli callouts` : Edit Obsidian callouts configuration
- `obs-cli archive` : Archive files in the vault

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

# Archive files
obs-cli archive
```

## License

MIT
