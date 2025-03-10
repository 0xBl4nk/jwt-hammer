# JWTHammer

A lightweight, CPU-based JWT signature brute force tool written in Go.

## Overview

JWTHammer provides an alternative to GPU-based tools like Hashcat for cracking JWT signatures. It's designed to be simple, efficient, and work on systems without dedicated graphics hardware.

## Features

- Pure CPU implementation - no GPU required
- Multithreaded design for better performance
- Support for HMAC-SHA256 signatures (HS256)
- Simple command-line interface
- Memory-efficient wordlist processing

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/jwthammer
cd jwthammer

# Build the binary
go build -o jwthammer

# Or run directly
go run main.go <jwt_token> <wordlist_file>
```

## Usage

```bash
./jwthammer <jwt_token> <wordlist_file>
```

### Example

```bash
./jwthammer "<full-token>" rockyou.txt
```

## Performance

Performance depends on your CPU and the size of your wordlist. JWTHammer automatically uses multiple workers to utilize available CPU cores.

## Disclaimer

This tool is intended for security testing and educational purposes only. Always obtain proper authorization before testing security on any system you don't own.
