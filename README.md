# MQTT HTTP Bridge

Extremely lightweight HTTP server for publishing MQTT messages via HTTP GET requests. Perfect for Raspberry Pi.

## Setup

1. Set environment variable:
```bash
export MQTT_BROKER=tcp://your-broker:1883
```

2. Build:
```bash
go build -o mqtt-bridge
```

3. Run:
```bash
./mqtt-bridge
```

## Usage

Trigger MQTT publish via HTTP:
```bash
curl "http://localhost:8080/publish?mqtt_topic=home/sensor&mqtt_message=hello"
```

## Cross-compile for Raspberry Pi

On your dev machine:
```bash
GOOS=linux GOARCH=arm GOARM=7 go build -o mqtt-bridge-pi
```

Transfer to Pi and run:
```bash
chmod +x mqtt-bridge-pi
MQTT_BROKER=tcp://broker:1883 ./mqtt-bridge-pi
```

## Resource Usage

Compiled binary is ~10MB, runs with <10MB RAM.
