package servers

import (
	"context"
	"fmt"
	"net/url"
)

func (c *client) CreateCryptokey(ctx context.Context, serverID, zoneID string) (*Cryptokey, error) {
	cryptokey := Cryptokey{}
	path := fmt.Sprintf("/api/v1/servers/%s/zones/%s/cryptokeys",
		url.PathEscape(serverID), url.PathEscape(zoneID))

	err := c.httpClient.Post(ctx, path, &cryptokey)
	if err != nil {
		return nil, err
	}

	return &cryptokey, nil
}
