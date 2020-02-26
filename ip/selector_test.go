package ip

import "testing"

func TestFilterIP(t *testing.T) {
	c := NewQingtingIP(
		"https://proxyapi.horocn.com/api/v2/proxies",
		"0KXT1659129263374359",
		"json",
		"unix",
		"no",
		"",
		20,
	)
	ips, err := c.Pull()
	if err != nil {
		t.Fatal(err)
	}
	r := FilterIP(ips)
	t.Log(r)
}
