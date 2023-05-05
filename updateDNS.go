package main

import (
	"fmt"
	"log"
	"os"

	client "github.com/asvany/hetzner-dns-go"
	"github.com/joho/godotenv"
)

func getZoneID(hDNSClient client.HetznerDNSClient, zoneName string) string {
	zones, err := hDNSClient.ZonesGet()
	if err != nil {
		log.Fatalf("error while retrieving zones list: %s\n", err)
	}

	for _, zone := range zones {
		if zone.Name == zoneName {
			return zone.ID
		}
	}

	return ""
}

func main() {

	err := godotenv.Load(".env", "common.env")

	ip := GetPodID()
	if ip == "" {
		log.Fatalln("no pod ip found")
	}
	fmt.Println("pod ip is: ", ip)

	if err != nil {
		log.Println("WARRNING: error while loading .env file: ", err)
	}

	HETZNER_DNS_TOKEN := os.Getenv("HETZNER_DNS_TOKEN")
	if HETZNER_DNS_TOKEN == "" {
		log.Fatalln("HETZNER_DNS_TOKEN environment variable is not set")
	}

	ZONE_NAME := os.Getenv("ZONE_NAME")
	if ZONE_NAME == "" {
		log.Fatalln("ZONE_NAME environment variable is not set")
	}

	hDNSClient := client.NewAuthApiTokenClient(HETZNER_DNS_TOKEN)

	zoneID := getZoneID(hDNSClient, ZONE_NAME)

	fmt.Println("zone id is: ", zoneID)
}
