# w1temp2mqtt

![GitHub license](https://img.shields.io/github/license/Noddy76/w1temp2mqtt.svg)
![Go](https://github.com/Noddy76/w1temp2mqtt/workflows/Go/badge.svg)

A simple daemon to report a 1wire DS18B20 temperature probe to MQTT.

## Usage

```
./w1temp2mqtt -broker tcp://mqtt.home:1883 -topic sensor/temperature -device 28-000000000000
```

| Option   | Default              | Description                    |
| -------- | -------------------- | ------------------------------ |
| device   |                      | 1wire device ID to report on   |
| interval | 10s                  | How often to read the device   |
| broker   | tcp://localhost:1883 | MQTT Broker connext string     |
| clientid | w1temp2mqtt          | MQTT client ID                 |
| topic    | w1temp2mqtt          | MQTT topic to publish lines on |
