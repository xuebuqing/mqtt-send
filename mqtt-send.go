package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/xuebuqing/ctrlapp/pkg/mqtt"
	"gopkg.in/yaml.v2"
)

/*
**read config.yaml
 */
func subCallBackFunc(Topic string, msg []byte) {
	fmt.Printf("Subscribe: Topic is [%s]; msg is [%s]\n", Topic, msg)
}

type econfig struct {
	Broker   string `yaml:"broker"`
	Clientid string `yaml:"clientid"`
	Passwd   string `yaml:"passwd"`
	Qos      byte   `yaml:"qos"`
	Usr      string `yaml:"usr"`
	Capath   string `yaml:capath`
	CertFile string `yaml:certFile`
	KeyFile  string `yaml:keyFile`
}

// create tls config
func CreateTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	pool := x509.NewCertPool()
	rootCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	ok := pool.AppendCertsFromPEM(rootCA)
	if !ok {
		return nil, fmt.Errorf("fail to load ca content")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientCAs:    pool,
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	}

	return tlsConfig, nil
}

func (c *econfig) getConf() *econfig {
	yamlFile, err := ioutil.ReadFile("conf/config.yaml")
	if err != nil {
		log.Printf("yamlfile.get err #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal :%v", err)
	}
	return c
}

//publish
func publish(t econfig, mes *mqtt.WillMessage, tls *tls.Config) {
	client := mqtt.NewMQTTClient(t.Broker, t.Clientid,
		t.Usr, t.Passwd, mes.Retained, tls)
retry_connect:
	err := client.Connect()
	if err != nil {
		time.Sleep(3 * time.Second)
		fmt.Println("recovering connecting")
		goto retry_connect
	}
	for i := 1; i < 10; i++ {
		fmt.Println("this is i:", i)
		client.Publish(mes)
		time.Sleep(time.Second)
	}

	client.Close()
}

//subscribe
func subscribe(t econfig, mes mqtt.WillMessage, tls *tls.Config) {
	// sub的用户名和密码
	client := mqtt.NewMQTTClient(t.Broker, t.Broker,
		t.Usr, t.Passwd, mes.Retained, tls)
retry_connect:
	err := client.Connect()
	if err != nil {
		time.Sleep(3 * time.Second)
		fmt.Println("recovering connecting")
		goto retry_connect
	}
	client.Subscribe(mes.Topic, mes.Qos, subCallBackFunc)

}
func main() {
	var t econfig
	t.getConf()

	tlsConfig, err := CreateTLSConfig(t.Capath,
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
	subscribe(t, *wmessage, tlsConfig)
	publish(t, wmessage, tlsConfig)

}
