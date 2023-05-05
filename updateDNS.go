package main

import (
    "fmt"
    "log"

    client "github.com/nl2go/hetzner-dns-go"
)

func main() {
    hDNSClient := client.NewAuthApiTokenClient("yourAuthAPIToken")

    zones, err := hDNSClient.ZonesGet()
    if err != nil {
        log.Fatalf("error while retrieving zones list: %s\n", err)
    }

    fmt.Println(zones)
}