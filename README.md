# NM File Decryptor

A simple Go tool to decrypt `.nm` configuration files used in NetMod Application. Supports multiple protocols such as `vless`, `trojan`, `ss`, etc.  

The tool reads `.nm` files from an input folder, decrypts their  encrypted configs, and writes the decrypted content into an output folder

---

## Features

- Supports any protocol prefixed with `nm-` (e.g., `nm-vless://`, `nm-trojan://`)
- Removes padding null bytes from decrypted output
- Saves decrypted files in a separate output directory

---

## Prerequisites

- Go 1.20+ installed
- Basic understanding of terminal commands  

---

## Installation

Clone the repository:

```bash
git clone https://github.com/<your-username>/nm-file-decryptor.git
cd nm-file-decryptor
```

# Usage
Command-line arguments
```
./nm-decryptor -input <input-folder> -output <output-folder>

Note: if you only have nm-xxx:// from clipboard (not the file), make one like "anything.nm" inside the input folder 
```


# Notes
- Input .nm files must be Base64 encoded
- Padding null bytes (\x00) are removed from the decrypted output automatically

License
This project is licensed under the MIT License ;)

Made by love :P
