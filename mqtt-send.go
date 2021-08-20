package main

import (
	"encoding/json"

	mqtt "github.com/xuebuqing/mqtt-send/pkg/mqtt"
)

func main() {
	var t mqtt.Econfig
	t.GetConf()

	tlsConfig, err := mqtt.CreateTLSConfig(t.Capath,
		t.CertFile, t.KeyFile)
	if err != nil {
		tlsConfig = nil
	}

	buf, _ := json.Marshal(t)

	wmessage := &mqtt.WillMessage{
		Topic:    "$hw/events/mqtt_test1",
		Payload:  buf,
		Qos:      0,
		Retained: false,
	}

	mqtt.Subscribe(t, *wmessage, tlsConfig)
	mqtt.Publish(t, wmessage, tlsConfig)

}
