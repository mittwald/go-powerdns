package zones

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mittwald/go-powerdns/pdnshttp"
)

type ZoneBasicDataUpdate struct {
	Kind       ZoneKind       `json:"kind,omitempty"`
	Masters    []string       `json:"masters,omitempty"`
	Account    string         `json:"account,omitempty"`
	SOAEdit    ZoneSOAEdit    `json:"soa_edit,omitempty"`
	SOAEditAPI ZoneSOAEditAPI `json:"soa_edit_api,omitempty"`
	APIRectify *bool          `json:"api_rectify,omitempty"`
	DNSSec     *bool          `json:"dnssec,omitempty"`
	NSec3Param string         `json:"nsec3param,omitempty"`
}

func (c *client) ModifyBasicZoneData(ctx context.Context, serverID string, zoneID string, update ZoneBasicDataUpdate) error {
	path := fmt.Sprintf("/servers/%s/zones/%s", url.PathEscape(serverID), url.PathEscape(zoneID))

	err := c.httpClient.Put(ctx, path, nil, pdnshttp.WithJSONRequestBody(&update))
	if err != nil {
		return err
	}
	return nil
}
