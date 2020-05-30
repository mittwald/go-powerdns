package servers

// Cryptokey model a PowerDNS server.
//
// More information: https://doc.powerdns.com/authoritative/http-api/cryptokey.html#cryptokey
type Cryptokey struct {
	ID         int      `json:"id"`
	Type       string   `json:"type"`
	KeyType    string   `json:"keytype"`
	Active     bool     `json:"active"`
	Published  bool     `json:"published"`
	DNSKey     string   `json:"dnskey"`
	DS         []string `json:"ds"`
	PrivateKey string   `json:"privatekey"`
	Algorithm  string   `json:"algorithm"`
	Bits       int      `json:"bits"`
}
