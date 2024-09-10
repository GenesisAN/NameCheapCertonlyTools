# Namecheap DNS Management Tool

This is a simple Go program for interacting with the Namecheap DNS API to update DNS records, specifically designed to manage DNS challenges for SSL certificates (e.g., `_acme-challenge`).

## Features
- Retrieve existing DNS records.
- Update DNS records dynamically.
- Simple flag-based input for API credentials and DNS information.

## Requirements
- Go 1.16+
- An active Namecheap account with access to API credentials.
- An internet connection to access Namecheap's API.

## Usage

### Flags

- `-u` : API User (required)
- `-k` : API Key (required)
- `-c` : Client IP (required)
- `-s` : Second-Level Domain (SLD, required)
- `-t` : Top-Level Domain (TLD, required)
- `-h` : Hostname (optional, defaults to `_acme-challenge`)
- `-v` : Value (required, the value to be set for the DNS record)

### Example Command

```bash
go run main.go -u <API_USER> -k <API_KEY> -c <CLIENT_IP> -s <SLD> -t <TLD> -v <VALUE> -h <HOSTNAME>
```

### Where:

`<API_USER>` is your Namecheap API user.  
`<API_KEY>` is your Namecheap API key.  
`<CLIENT_IP>` is your public IP address.  
`<SLD>` is the Second-Level Domain (e.g., example for example.com).  
`<TLD>` is the Top-Level Domain (e.g., com for example.com).  
`<VALUE>` is the DNS record value to set (e.g., a challenge value).  
`<HOSTNAME>` is the optional hostname to update (e.g., _acme-challenge).  
### Example
```bash
go run main.go -u myUser -k myAPIKey -c 192.168.1.100 -s example -t com -v "new_value" -h "_acme-challenge"
```
