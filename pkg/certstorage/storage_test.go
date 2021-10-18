package certstorage

import "testing"

func TestAcmeCertificate_AsPem(t *testing.T) {
	type fields struct {
		Domain            string
		CertURL           string
		CertStableURL     string
		PrivateKey        []byte
		Certificate       []byte
		IssuerCertificate []byte
		CSR               []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Domain:        "domain",
				CertURL:       "certUrl",
				CertStableURL: "stableUrl",
				PrivateKey: []byte(`private key
so much private
wow private`),
				Certificate: []byte(`this is me
i totally swear`),
				IssuerCertificate: []byte(`so much responsibility
awesome issuer`),
				CSR: nil,
			},
			want: `this is me
i totally swear
private key
so much private
wow private`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cert := &AcmeCertificate{
				Domain:            tt.fields.Domain,
				CertURL:           tt.fields.CertURL,
				CertStableURL:     tt.fields.CertStableURL,
				PrivateKey:        tt.fields.PrivateKey,
				Certificate:       tt.fields.Certificate,
				IssuerCertificate: tt.fields.IssuerCertificate,
				CSR:               tt.fields.CSR,
			}
			if got := cert.AsPem(); got != tt.want {
				t.Errorf("AsPem() = %v, want %v", got, tt.want)
			}
		})
	}
}
