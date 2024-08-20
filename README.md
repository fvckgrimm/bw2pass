# Bitwarden to Password Store Converter (bw2pass)

bw2pass is a Go tool that converts a Bitwarden JSON export into a Password Store (pass) compatible format. It organizes your passwords and secure notes into a hierarchical structure based on folders, domains, and entry types.

## Features

- Converts Bitwarden login entries and secure notes to Password Store format
- Organizes entries by folder structure (if present in Bitwarden)
- Creates a domain-based hierarchy for login entries
- Stores secure notes in a separate 'notes' directory
- Handles multiple entries with the same name by appending a counter
- Sanitizes folder and entry names for filesystem compatibility
- Preserves additional information such as URIs, TOTP, and notes

## Prerequisites

- [Go](https://go.dev/dl/)
- [pass](https://www.passwordstore.org/) (Password Store) installed and initialized on your system

## Installation

### Pre-Built Binaries

Pre-Built binaries can be found under [releases](https://github.com/fvckgrimm/bw2pass/releases)

### Building 

1. Clone this repository:


```bash
git clone https://github.com/fvckgrimm/bw2pass.git
cd bw2pass
```

2. Build the tool:


```bash
go build -v -o bw2pass main.go
```

## Usage

1. Export your Bitwarden vault as a JSON file.

2. Run the tool with your Bitwarden JSON export file:

```bash
./bw2pass path/to/your/bitwarden_export.json
```

3. The tool will process your Bitwarden export and insert entries into your Password Store.

## Output Structure

The tool organizes your Password Store as follows:

```
- Login entries:
Password Store
├── folder_name_1
│   ├── example.com
│   │   └── login_entry
│   └── another-example.com
│       └── another_login_entry
└── example.com
└── login_entry

- Secure notes:
Password Store
├── folder_name_1
│   └── notes
│       └── secure_note_1
└── notes
├── secure_note_2
└── secure_note_3
```

## Notes

- The tool uses the `pass` command-line utility to insert entries. Make sure you have initialized your Password Store before running this tool.
- Entries with the same name in the same location will have a counter appended (e.g., `entry_name_1`, `entry_name_2`).
- Folder and entry names are sanitized to ensure filesystem compatibility.
