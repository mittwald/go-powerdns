package zones

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mittwald/go-powerdns/pdnshttp"
)

func (c *client) AddRecordSetToZone(ctx context.Context, serverID string, zoneID string, set ResourceRecordSet) error {
	path := fmt.Sprintf("/servers/%s/zones/%s", url.PathEscape(serverID), url.PathEscape(zoneID))

	set.ChangeType = ChangeTypeReplace
	patch := Zone{
		ResourceRecordSets: []ResourceRecordSet{
			set,
		},
	}

	return c.httpClient.Patch(ctx, path, nil, pdnshttp.WithJSONRequestBody(&patch))
}

func (c *client) AddRecordSetsToZone(ctx context.Context, serverID string, zoneID string, sets []ResourceRecordSet) error {
	path := fmt.Sprintf("/servers/%s/zones/%s", url.PathEscape(serverID), url.PathEscape(zoneID))

	for idx := range sets {
		sets[idx].ChangeType = ChangeTypeReplace
	}
	patch := Zone{
		ResourceRecordSets: sets,
	}

	return c.httpClient.Patch(ctx, path, nil, pdnshttp.WithJSONRequestBody(&patch))
}
