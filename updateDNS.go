package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func checkPort(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Printf("Can't reach %s:%d : %v\n", host, port, err)
		return false
	}
	defer conn.Close()
	fmt.Printf("Connection successfull %s:%d!\n", host, port)
	return true
}

func main() {

	err := godotenv.Load("secret.env", "unsecret.env")
	if err != nil {
		log.Println("WARRNING: error while loading all env files: ", err)
	}

	HETZNER_DNS_TOKEN := os.Getenv("HETZNER_DNS_TOKEN")
	if HETZNER_DNS_TOKEN == "" {
		log.Fatalln("HETZNER_DNS_TOKEN environment variable is not set")
	}

	ZONE_NAME := os.Getenv("ZONE_NAME")
	if ZONE_NAME == "" {
		log.Fatalln("ZONE_NAME environment variable is not set")
	}

	HOST_NAME := os.Getenv("HOST_NAME")
	if HOST_NAME == "" {
		log.Fatalln("HOST_NAME environment variable is not set")
	}

	RECORD_TYPE := os.Getenv("RECORD_TYPE")
	if RECORD_TYPE == "" {
		RECORD_TYPE = "A"
		log.Println("RECORD_TYPE environment variable is not set, defaulting to A")
	}

	ip, status := GetIpAndConnectionStatus()
	if !status {
		log.Fatalln("service not ready, exiting")
	}

	hetznerClient := NewHetznerClient(HETZNER_DNS_TOKEN, ZONE_NAME, HOST_NAME, ip).findZoneID()

	log.Println("Zone ID is: ", hetznerClient.zone_id)

	// log.Println("setup ip is: ", ip)
	record_id, record_ip, record_type := hetznerClient.findRecordDataByName(hetznerClient.host_name)
	log.Println("record_id: ", record_id)
	log.Println("record_ip: ", record_ip)
	log.Println("record type: ", record_type)
	if record_ip == ip {
		log.Println("ip is the same, exiting")
		return
	} else if record_ip == "" {
		log.Println("record ip is empty, creating new record")
		hetznerClient.createRecord(hetznerClient.host_name, ip, RECORD_TYPE)
	} else {
		log.Println("record ip is different, updating record")
		hetznerClient.updateRecord(record_id, hetznerClient.host_name, ip, record_type)
	}

}

func GetIpAndConnectionStatus() (string, bool) {
	ip := GetPodID()
	if ip == "" {
		log.Fatalln("no pod ip found")
	}
	log.Println("pod ip is: ", ip)

	timeoutSecs := os.Getenv("PORT_CHECK_TIMEOUT")
	if timeoutSecs == "" {
		timeoutSecs = "3"
	}
	timeoutSecsFloat, err := strconv.ParseFloat(timeoutSecs, 64)
	if err != nil {
		log.Fatalln("Error while converting string to float:", err)

	}
	timeout := time.Duration(time.Duration(timeoutSecsFloat)) * time.Second

	portStr := os.Getenv("CHECK_PORT")

	if portStr == "" {
		portStr = "443"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalln("Error while converting port to int:", err)
	}

	if checkPort(ip, port, timeout) {
		fmt.Println("The port is open.")
		return ip, true
	} else {
		fmt.Println("The port is closed.")
		return ip, false
	}
	// return "", false
}
