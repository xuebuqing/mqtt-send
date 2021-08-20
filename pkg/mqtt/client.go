package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Econfig struct {
	Broker   string `yaml:"broker"`
	Clientid string `yaml:"clientid"`
	Passwd   string `yaml:"passwd"`
	Qos      byte   `yaml:"qos"`
	Usr      string `yaml:"usr"`
	Capath   string `yaml:capath`
	CertFile string `yaml:certFile`
	KeyFile  string `yaml:keyFile`
}

func subCallBackFunc(Topic string, msg []byte) {
	fmt.Printf("Subscribe: Topic is [%s]; msg is [%s]\n", Topic, msg)
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

func (c *Econfig) GetConf() *Econfig {

	//Obtain the config.yaml directory
	filePath, _ := os.Getwd()
	filePath = filePath + "/conf/config.yaml"

	yamlFile, err := ioutil.ReadFile(filePath)
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
func Publish(t Econfig, mes *WillMessage, tls *tls.Config) {
	client := NewMQTTClient(t.Broker, t.Clientid,
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
func Subscribe(t Econfig, mes WillMessage, tls *tls.Config) {
	// sub的用户名和密码
	client := NewMQTTClient(t.Broker, t.Broker,
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
