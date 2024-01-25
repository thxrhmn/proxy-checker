# PROXY CHECKER

[![asciinema CLI
demo](proxy-checker.png)](https://asciinema.org/a/s0eEzhOAeFpDNqSHSkb1by3jh)

This project is a simple tool to check the status of proxies. It helps you determine whether a proxy is still usable or no longer functional.

## Key Features
1. **Proxy Availability Check**: This project can verify whether a proxy is still active or has become unusable.
2. **Goroutine Usage**: Implementation utilizes goroutines to check multiple proxies simultaneously, enhancing efficiency and speed of the checking process.

## System Requirements
1. Golang version 1.21 or later.
2. Stable internet connection to perform online proxy checks.

## How to Use
1. **Clone the Project**:
```bash
git clone https://github.com/thxrhmn/proxy-checker.git
```
2. **Navigate to the Project Directory**:
```bash
cd proxy-checker
```

3. **Configure Proxies**:
Edit the proxy_list.txt file and add the list of proxies you want to check. Each line should contain one proxy, and its port should be separated by a colon (:).
```bash
127.0.0.1:8080
192.168.1.1:3128
```
4. **Download packages**:
```bash
go mod tidy
```

5. **Run the application**:
```bash
go run main.go -p proxylist.txt

```
The application will start checking the proxies and provide results indicating whether each proxy is still usable.