package servers

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *client) GetCryptokey(ctx context.Context, serverID, zoneID string, cryptokeyID int) (*Cryptokey, error) {
	cryptokey := Cryptokey{}
	path := fmt.Sprintf("/api/v1/servers/%s/zones/%s/cryptokeys/%s",
		url.PathEscape(serverID), url.PathEscape(zoneID), url.PathEscape(strconv.Itoa(cryptokeyID)))

	err := c.httpClient.Get(ctx, path, &cryptokey)
	if err != nil {
		return nil, err
	}

	return &cryptokey, nil
}
