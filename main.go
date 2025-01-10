package main

import (
	"bufio"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

func getIpAddr() (string, error) {
	cfEndpoint := "https://cloudflare.com/cdn-cgi/trace"

	req, err := http.NewRequest("GET", cfEndpoint, nil)

	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error while closing response body", err)
		}
	}(resp.Body)

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ip=") {
			return strings.TrimPrefix(line, "ip="), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("could not find ip address in response")

}

func getRootDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return domain
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfToken := os.Getenv("CF_TOKEN")
	if cfToken == "" {
		logger.Error("CF_TOKEN is required")
		os.Exit(1)
	}
	recordName := os.Getenv("RECORD_NAME")
	if recordName == "" {
		logger.Error("RECORD_NAME is required")
		os.Exit(1)
	}

	proxied := os.Getenv("PROXIED")
	proxy := false
	if proxied == "" || proxied == "0" {
		proxy = false
	} else {
		proxy = true
	}

	interval := os.Getenv("INTERVAL")

	if interval == "" {
		interval = "5m"
	}

	intervalTime, err := time.ParseDuration(interval)
	if err != nil {
		logger.Error("INTERVAL is not a valid duration", err)
		os.Exit(1)
	}

	logger.Debug("Starting cloudflare api session")
	api, err := cloudflare.NewWithAPIToken(cfToken)

	if err != nil {
		logger.Error("Error creating cloudflare api client", err)
	}

	logger.Info("Getting zone id")
	zoneId, err := api.ZoneIDByName(getRootDomain(recordName))
	logger.Info("zone id: ", zoneId)

	if err != nil {
		logger.Error("Error getting zone id", err)
	}

	cfZone := cloudflare.ZoneIdentifier(zoneId)
	logger.Info("begin updating dns record")
	for {
		logger.Info("Getting ip address")
		ip, err := getIpAddr()

		if err != nil {
			logger.Error("Error while getting external IP", err)
			continue
		}
		logger.Info("External ip address: ", ip)

		ctx := context.Background()

		logger.Info("Searching DNS records")
		records, _, err := api.ListDNSRecords(ctx, cfZone, cloudflare.ListDNSRecordsParams{Content: ip, Type: "A", Name: recordName})
		if err != nil {
			logger.Error("Unable to search for dns record", err)
			continue
		}
		if len(records) == 0 {
			logger.Info("No records found for: ", ip)
			logger.Info("Creating new record")
			newRecord := cloudflare.CreateDNSRecordParams{
				Content: ip,
				Name:    recordName,
				Type:    "A",
				Proxied: &proxy,
			}

			logger.Info("Record creation info: ", newRecord)
			logger.Info("Creating new record")
			_, err = api.CreateDNSRecord(ctx, cfZone, newRecord)
			if err != nil {
				logger.Error("Error while creating record", err)
				continue
			}
			continue
		}

		record := records[0]

		logger.Info("Found record: ", record)

		if record.Content == ip {
			logger.Info("IP address has not changed")
			<-time.After(intervalTime)
			continue
		}

		logger.Info("IP address has changed")
		logger.Info("Updating DNS record")
		updatedRecord := cloudflare.UpdateDNSRecordParams{
			Content: ip,
			Name:    recordName,
			Type:    record.Type,
			Proxied: record.Proxied,
			TTL:     record.TTL,
			Tags:    record.Tags,
		}
		logger.Debug("Updated record: ", updatedRecord)

		_, err = api.UpdateDNSRecord(ctx, cfZone, updatedRecord)

		if err != nil {
			logger.Error("Error while updating the record", err)
			continue
		}

		logger.Info("Updated DNS record to", ip)

		// sleep for 5 minutes
		<-time.After(intervalTime)
	}

}
