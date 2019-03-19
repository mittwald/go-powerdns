package zones

import (
	"context"
	"fmt"
	"net/url"
)

func (c *client) GetZone(ctx context.Context, serverID, zoneID string) (*Zone, error) {
	zone := Zone{}
	path := fmt.Sprintf("/api/v1/servers/%s/zones/%s", url.PathEscape(serverID), url.PathEscape(zoneID))

	err := c.httpClient.Get(ctx, path, &zone)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}
