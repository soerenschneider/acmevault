package acme

import (
	"reflect"
	"testing"
)

func Test_fixLineBreaks(t *testing.T) {
	cert := []byte(`-----BEGIN CERTIFICATE-----
example data
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
more example data
-----END CERTIFICATE-----

`)
	wanted := []byte(`-----BEGIN CERTIFICATE-----
example data
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
more example data
-----END CERTIFICATE-----
`)

	got := fixLineBreaks(cert)
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("Expected %s, got %s", string(wanted), string(got))
	}
}
