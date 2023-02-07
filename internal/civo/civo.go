package civo

import (
	"os"

	"github.com/civo/civogo"
	"github.com/rs/zerolog/log"
)

// GetDNSInfo try to reach the provided hosted zone
func GetDNSInfo(domainName, region string) (string, error) {

	log.Info().Msg("GetDNSInfo (working...)")

	client, err := civogo.NewClient(os.Getenv("CIVO_TOKEN"), region)
	if err != nil {
		log.Info().Msg(err.Error())
		return "", err
	}

	civoDNSDomain, err := client.FindDNSDomain(domainName)
	if err != nil {
		log.Info().Msg(err.Error())
		return "", err
	}

	return civoDNSDomain.ID, nil

}
