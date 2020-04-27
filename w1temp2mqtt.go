/*
Copyright 2020 James Grant

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	var deviceID = flag.String("device", "", "1wire device ID to report on")
	var interval = flag.Duration("interval", 10*time.Second, "How often to read the device")

	var mqttBroker = flag.String("broker", "tcp://localhost:1883", "MQTT Broker connext string")
	var mqttClientID = flag.String("clientid", "w1temp2mqtt", "MQTT client ID")
	var mqttTopic = flag.String("topic", "w1temp2mqtt", "MQTT topic to publish lines on")

	flag.Parse()

	if len(*deviceID) == 0 {
		log.Fatal("device must be specified")
	}

	if strings.Index(*deviceID, "28-") != 0 {
		log.Fatalf("Device must be a family 28 device (DS18B20 or similar)")
	}

	opts := mqtt.NewClientOptions().
		AddBroker(*mqttBroker).
		SetClientID(*mqttClientID)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer mqttClient.Disconnect(250)

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for range ticker.C {
		data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + *deviceID + "/w1_slave")
		if err != nil {
			log.Printf("Problem reading value from device (%v)", err)
			continue
		}

		raw := string(data)

		var i = strings.Index(raw, "t=")
		if i == -1 {
			log.Println("Unable to extract temperature from device")
			continue
		}
		raw = strings.Trim(raw[i+2:], "\t\r\n ")

		temperature, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			log.Printf("Unable to parse temperature value (%v)", err)
			continue
		}
		temperature = temperature / 1000.0

		text := fmt.Sprintf("{\"temperature\":%0.3f}", temperature)
		log.Printf("%s : %s", *mqttTopic, text)
		token := mqttClient.Publish(*mqttTopic, 0, false, text)
		token.Wait()
	}
}
