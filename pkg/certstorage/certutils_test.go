package certstorage

import (
	"encoding/base64"
	"reflect"
	"testing"
)

const dummyRsaPrivate = `MIIEpAIBAAKCAQEAqv2BPRLptl6zrPU+GnjWZ5s+6A+57Rlc6/TR2KIhz/cVyRT3
bkbRglcM6LGMwLGx5tlqXHWpRgJf32tKtv2CuKjn/Jn2h3qpyTjIvbLgg6/cNph1
m+5OVXxBU5aMAl1p5r4RqNWKrERbvHCDlLi+l1M6hEgrDEpIw6I8Cs5L5Q/JQh+U
C1M3kIBgk0G6Ny/xriE5d234zuH6H7Br5VjtcnpKHsuxVqtn6LX0nBjZQtmoO60P
6CbyRUFS5alMIdhCkUoZAe94+iI4gBUb5ukNotvqyKzJOEC6+27jJ6jER/lGJReB
HnQmqQw7/IRyQXPwnFIIRoKe5OmWj+Sgl985IwIDAQABAoIBAETxzH+MtbG0A6yU
ggL4gDLsVPQLC0O/u0jkwQwha5LIJP+cNZxAb8+nO+xDUuaLf5j+RzWat7Mj1/Zk
845PL7s3V7rxdYbw/a7F96MNkhtm+FsHJDzIXMt5O3nxtOhrM/023DYATFbjhT24
/EUmLKslgu29j6X3+djv/Fe7ELz99MDQh+gEp4wJ8M3w9J/+nanKcDQv5EkPM1f2
S8TfyLEaTO4QC0/QwRWRJGe0OYP3uTbo/Jz+nCVOEN51ti/VWoKGq6e8521zLHds
4LxPFg9m/OSWAbzSGbUM4F1Ju+cGvtlWsnCttGSBxtEj6ASoe3ujYPNSCGAT7sKr
12zlD2ECgYEA1LYoV+2YQL0qfA0eAjPdSwYpd6z015/xSGOYZ27ABx1EiJYi+Xcr
mg88oSWo1K7BfvlXq76RVT8i82bqpVEailygUJuhUDZPGgXu+nkrGbkNTejEaelo
TN+b6WLXn0MM9P7GML56rl7lghqScXsj8JFqruZl5VHLmTkcdSTJnxMCgYEAzcnA
c/r3D9A2GU3ig7MJqp1PbAJVAkoZnw4awILsqiTH4G8zJL7wGyK4nd3M3RRZSlsu
SUdrsWNy5bleoS53HJAQ2XmqcY48VE5xp7MgsmjNzkvUqggY1tKWotQ1vaccIxzK
C6BtDeVvcIcEsu2BJR3auGB7Am0DFpHhnl5lb7ECgYB9H7X1mx6/nkbaeADZ/NZE
pphH/KZ+HGibU3K4DS7KQI4q5l4mbtJoKmRYysVYboAGB/hpX40wHxaZJUYL/vxk
vX3UTOgEPh4jya+9CP+wfulzlePMBW/EIJkfHXSHC7MYfrHoxHN4FPxenvDb7jrC
7VdbEM6JqabhN/YkdoJfnwKBgQDDzx/fe8IP+ucvBqvs8rPU8yA6PrHSnt1oOcx8
t1cnwh2p0UVRxpjCyTPvire++PjJRp9xPJMdU+pk6hA/v+24cOCHfjwPFu8SrR16
G3iWKiECFaeCLngsGG8a+l80WUjtpBwgGgMKgPCbmu1+r3z96z1NyQfi1AVnOWH8
Bgfw0QKBgQCjy1Ui5OgOrACYppuQD6lMzB2OGRfFABbA4bBCGB6veWrY+qToTlMp
vPvNfw37WF9h3wdMz8ObX6FD2iXWd7fFgaPkyldWg4Val8R/e0FsGjCFWeF2goZG
F/WDPCNexoH9oQBKHM6zFEk8XEi11aIfN6Blmv+lTNWKRqdZfg9srA==`

const dummyCert = `MIIEpAIBAAKCAQEAuF5FvRNedBjT8rZ8Vy69+tRwidtR3e8g2mFqZCmuL3Cq+ps3
DR4K7q51nmCTQhUOOTKW3z+UbaQ0cJnfPoIkDm71x2A3SP4EafBhi19GUt04stpK
eG2TjvOz/9wWPD0ydK/Wn8Ktm5MJVRYjgnX4D6gSgvH39z+7xjYczr7n/rlyu7EN
I42Y11bxndmlWFFN7TFr8ps4MM3FAb1gTqt8JeeDYbimDvCsP15XmVS32OSV85yT
45GXn4lxUDcV+UX3lefdmfdrnxtac4WNex4gKYutinlrhA4krT0SmYkv8AVJABX0
T1omU66Hpm4px+PG2OlRI+5dPKoTNdzgKE520QIDAQABAoIBAQCztxlNrDK3Anif
r6hWp3aCqVAl4QbVSkKA/NJjXomvLqlL7b7k53MKpD58qnEAylt5/8k3RtedZKHF
XQhS+zLAawLjXPOnEA0nYy0CWVXIcmHRXZ2X9GAQyZedAWEfSRwEyF5yjqcB7nr7
WMA2bF3ojRG2WL04YwTbumWsVkT4cw8vEbV+zZcMvssiwtd9fD00wgPsicilrKa1
m5cZI2HmATtsHO1UMz51sFmtwZk6N6DmmuSrjJp54qvdTCPfMqzLeG0sEtEToApb
lk+X9ZDRqfJSq7BYYnSF/984cEaSwBDg6PZVvjoLsSMfzAr9LTgO1o4RRuEUycTN
zxVT7TE9AoGBAN5KHnvvu8hwPT0RmBukLLP2T3dNZUKw/hG2V6Lf4uSvrJoqAZgk
mjXFEemAiQyXxQB3diUgaFcZAEe9YlS7xiAczpGMvY4f4hQl3wPrEobXWlVU+Le5
ReyJB1xPXOw5lNfoE+7BXmlwbg4/A8SXHTK4cONaGpPcwL6D0TVwR7kbAoGBANRT
8mFOB/pJNNXYtjz9eV+mWLXDoa3d0Wfsd3g1EGowhKXpyH8ZCGbF9i3NinhHJqX7
cI0SbdAzxZsb9zuy4TlHfcIW5lcFd/s4NGRXb+tQZPeMwRNwTS7IUwI9RXZnvJ5X
MSmzrpkAnu159wBNHFSYLs152nIDJHf1KB5LShqDAoGBAKAL2Zr+ZsPQHCdmDZCJ
Si5gf6K0NN4ptSRnlv2MGT7yRWHSMMesQuPH+jeW9hX9CDhnysa8aKOdOphsgOnl
MKdaVlhxbuFvj0VWZxXZIMO5Ni8OOO1/FhtSJdyOv6bNZp91VzSmQSXFb1gOgtX+
v2UPaygmbdBcNuJ04iKJrvpVAoGAR0AETfLJaou2Vyxpuv20BQzlJ9mGH7oX0y/e
x0/HOlsdqC9TQJM95n73pVcb6FC3/2ro0e4lO4CkBvDTfg0A1x/Oa5sToJANOxgZ
PLK6s9Q+jXOGNuewfspqUI4PCTS2bswDi2LobB9xNW+AG3HE1/5Zdko1q5yyWC7E
T6YkL9ECgYAVC+Ui9Ag3LsSfUs5GbMpPNH4XPYpHZcE/bXarqjvVb5d9d5tXsuK0
4gVjAHsqStkiVouwuK5h4h5sIiJMA6E+fWnDZHC6Xx3S6FHjR62/V7s96ZNG2ycy
u+H1ib0m9OqCB/Pw0omI+WHOngrMmLhos9nDS6TgXuwVTpOoep3Zog==`

const dummyIssuer = `MIIEpAIBAAKCAQEAvmq0+10V1sMF6k8J4I21n3uWjU9T+KdiDRHb/JgzOrx6uS+v
CQDksAPHMnH4ue6gsCneRtkuhUyitEZv2gwVH3noNnZWNfxs8z1AaX37Fa/Uk1CP
OwFUh2dsFQgPe0JOmvwEfhlETZyHI+TelKR1XDXw531JOyt5FAsKAg3wBlLkHkyj
sul37srD1OJE/fjqBg3eWultUHTwCD59GtnbZqDlL4dQjMMchN1b0nil2WjVBPGX
ovxhlKJ8z0uQpWmICqY0c709fV49Pvcd2t80gs2O7Hf/7zEOucvXajK/7LrzU+dO
mPvZI8Dyx028j4f8vqB13L/QJLSoyCz/8jVBWQIDAQABAoIBAQCai7K+POvHtdus
M2A53+okOcOUh2kI7Jl5MCCTH9iceHNGsDvpG8+ASGC5QaV1CwdiU2jzqbvHNs7r
cCvCFoJiXKgq49rO0ESBGxqXREewb1giBIVrh4XarAcd/r/J86QmyBrBWbKFJ2DL
sisxC1WxdJpE1/vCyWLo8Ji72CISjnbXEOiBXpvP7iwE0tbvJ92U7TJcSBlJz/IQ
UFohvVDb/iC420mWJkKbkPtQE58J+XzWsvFEjODSYM33BLKm4siSkjGpYpKZV4OI
J7k00j7PdvsvrnaYSAjGkZddujPvisjGjP8xLyaOn3mvOKBEheFgKGKQvbzg8JEF
ovWm3XkBAoGBAO7E8BfuXY5XW29wTaZo5Qn9Cz72wPcAa3IFAu/SUAnS98QdR7Ln
10TAHXc8EhUxpxjtKZlMFFRvi2WaeZytsvFynzU3lsTSBzca3woaNV2A3riQ78kg
Q2kS0bbtfEXi/BOhVna6JATS3qNGzH+mc+GPxTRVTDzlGYOiwRWybUUhAoGBAMwo
f3cYYr9RVzexpP02ZnvrXnb5uX/RhU2h9onSjdvLAO0sqU3A87aIKaQtSxLTT3vI
Y7zAomI17qCeYv1ZkN4ngkvQbFqb0qIEI+3wtTgeIeYCssNb6alpuh8G/10yJ9Sv
IqUQspuqQOpHmeCZccsiuDmALEjLmZOIWEb5aD05AoGBAJIrDbAYxD03TUpTPbX5
0PzkY9YPyOFs6FnMp5eY8FaTSApOwm3LcAUudttfctJ2qwyfYy/tWyS1hGiWwIwh
6cHVoZE6jpm+2ZvqX1AX60NqeO4UDDbcAWh5lNifWcyOwDJkOkJEgXhSfukFlnsu
sxIKqXb4IMvGlG/5WqqlqC8BAoGAUlibYSwi2Ew0w7ARflis6Zq8FX0Qhy+5duC3
EkwtD9RH4WI8P7JuGte9BA2I1GULEEB5ii6g0MA0KfD4uHuh5RlGgtHkgn+La/ID
k/uc/K+auK2p8QZnrv+IJO+rnKmYSz8A2Fdt0z/OwzByLpd1wJuWdwrt0cbdgRZj
lO0QUHECgYAPdZuMSgZkIiUlQ1zPzxo0KBJtKkaKTktJ57R7TgAHSzIddk6DnnZF
C9cu+isNLQkPO9VBQsag88ANsIO9cdLNLzCXYTf/gpo/8UmJWmFwrpKQinuQo5Xb
mh2ZO6+zpTijFdJJ5f6nxy8mU/+iQ0Yc+B0GORnNwwb5PhiFD5d2kA==`

const dummyCsr = `MIIEowIBAAKCAQEA5PrvLkhSK5hxHpv1f4s2NH4SdMZZbMPBX31POm1nQG7X5CZ0
d4cr/mRDM/kOGVNpC3TPt5Bnl8H31DwmgKYIw2NbVy7YIlhieW2yUIi4NDljat7i
cr1Sm3+w91TT5+qgfh7kNOV4p84nTovyVIEsYEUbEvbYhD/lFirOskHDEhAY73PE
a6JszeGjCuWMr8HapcG5S1075codevk8nR+GUFzc7GaLJKbH3scFm6Uf97A7CMlz
n0avppkAxKuU5/+lYpgjQ69SZbebnuAlaXvK0Jm3BwugHHSCnPZ/FoSCh2xsUvlf
qObQ9SxODFYC7P1SyFuXOIhHuJy3S83P5ErAiQIDAQABAoIBACdrE1WyWYLrwT11
t7N3MaOjuGWl56sTn+xiVAtI3id3bW73N8GD4YzvkaoWy9iMRV8VgtSk5VB4scM1
f6NR8dxA9G5zv/1Zncotmi7G+n7zCixRpkX+VYQzXTGWxsv71hkgFEUO49BvatDY
wqTNf+gqvVhsaiWKIlPIlsCVFZG8JEUCn3xb/lNmEfhDdvn6H4U2eYVh1Dlxoies
Cey1+otXs7P/W3M6KWX6WMR9tM85SGe6Svpif6peFq3x0C3BLBER560xBgq+abTy
XtaYPWfDpiN0O84G5w8MWe4z2FNkIAfbF7wA2zL5YkBpgt49NEv/ADxdh4iIRtSY
IQLU2MkCgYEA+vQshNHij3ThbzOcEUomRlf/GOEcgYQVdj5c/zXrAIaZWfGpKLpZ
/+1qNSk9TgW6Q52xTEjFdXXBGPkuYARLxHhIVBbTLY/hwLjjQlxrXcyyAwf+Z/Tw
3Yxyn8GN/OjMJUqUvIs0noBvjwNGNAgl6/wrpDIV68PUaAD8XpTmNXsCgYEA6ZWl
0DX3We9jSSm2sDv3q7xm2PbccVVGM9iVfLGt5Bt9vFlh88K2gcdvcwgcLyOEEOZt
RjpJ/wnSBwRXcio9tTT3QIPAKRtCcawIgnTIcy/m19K0eAvcS+9x7J3CxkLO2kCk
Lmn4SPgmgQC3gtbcsvd/aHFfQXOhKNV73G9viMsCgYEArlH2qrxwuF3kSq569ref
JrXxiYK1nnH1xpFDYDQ/7bmRxJzNeHEaG2D7qbnfz9bCsD1V7zuNji4h2AsfX6sc
RnmXJHJGdxu/IXRMyMgR/LI35UskOWo39m2dIcP9sXS0eiL4do/sGT32QE1x8qrG
TMp6NjBkccUyQpyMsdaUowkCgYAKAS3b8CPLB1TSUmYPwFHIWkZxbolclVFvcQxe
DeIrzf2hrpZicWmNv3QHkkZawoOqkaQGiQKYWNxVDpuMOUDxXPZmHf6CBDfhVIP8
ynG3dUrG3bB7H87stbHEd7Fa+ouPj4s4rbNDtNU5W0WA5iEHzU/4sjppPEGf1Rz8
AQ3e5wKBgArwE9Z+LIMY4GirDyBKHHthBQ5a1wlndDNot/r3LMv5d9imeZl8qteH
YPOVK2/hqXzN+3pdq+vJFzMS0MLogynzDwCDWib1jH9HfFPDMAZjp7urMHKzgoZM
mtyMdGDPrwtHw0sGs6rFlk6DC6YBKAvvFssnIyRyWy3GhLppymul`

func TestMapToCert(t *testing.T) {
	decodedPrivateKey, _ := base64.StdEncoding.DecodeString(dummyRsaPrivate)
	decodedCert, _ := base64.StdEncoding.DecodeString(dummyCert)
	decodedIssuer, _ := base64.StdEncoding.DecodeString(dummyIssuer)
	decodedCsr, _ := base64.StdEncoding.DecodeString(dummyCsr)

	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *AcmeCertificate
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil",
			args:    args{nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "incomplete",
			args: args{map[string]interface{}{
				"bla": "bla",
			}},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{
				map[string]interface{}{
					vaultCertKeyCsr:        dummyCsr,
					vaultCertKeyIssuer:     dummyIssuer,
					vaultCertKeyUrl:        "url",
					vaultCertKeyCert:       dummyCert,
					vaultCertKeyDomain:     "domain.tld",
					vaultCertKeyPrivateKey: dummyRsaPrivate,
					vaultCertKeyStableUrl:  "stable url",
				},
			},
			want: &AcmeCertificate{
				Domain:            "domain.tld",
				CertURL:           "url",
				CertStableURL:     "stable url",
				PrivateKey:        []byte(decodedPrivateKey),
				Certificate:       []byte(decodedCert),
				IssuerCertificate: []byte(decodedIssuer),
				CSR:               []byte(decodedCsr),
			},
			wantErr: false,
		},
		{
			name: "invalid csr",
			args: args{
				map[string]interface{}{
					vaultCertKeyCsr:        "invalid",
					vaultCertKeyIssuer:     dummyIssuer,
					vaultCertKeyUrl:        "url",
					vaultCertKeyCert:       dummyCert,
					vaultCertKeyDomain:     "domain.tld",
					vaultCertKeyPrivateKey: dummyRsaPrivate,
					vaultCertKeyStableUrl:  "stable url",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid issuer",
			args: args{
				map[string]interface{}{
					vaultCertKeyCsr:        dummyCsr,
					vaultCertKeyIssuer:     "invalid",
					vaultCertKeyUrl:        "url",
					vaultCertKeyCert:       dummyCert,
					vaultCertKeyDomain:     "domain.tld",
					vaultCertKeyPrivateKey: dummyRsaPrivate,
					vaultCertKeyStableUrl:  "stable url",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid cert",
			args: args{
				map[string]interface{}{
					vaultCertKeyCsr:        dummyCsr,
					vaultCertKeyIssuer:     dummyIssuer,
					vaultCertKeyUrl:        "url",
					vaultCertKeyCert:       "invalid",
					vaultCertKeyDomain:     "domain.tld",
					vaultCertKeyPrivateKey: dummyRsaPrivate,
					vaultCertKeyStableUrl:  "stable url",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid private key",
			args: args{
				map[string]interface{}{
					vaultCertKeyCsr:        dummyCsr,
					vaultCertKeyIssuer:     dummyIssuer,
					vaultCertKeyUrl:        "url",
					vaultCertKeyCert:       dummyCert,
					vaultCertKeyDomain:     "domain.tld",
					vaultCertKeyPrivateKey: "invalid",
					vaultCertKeyStableUrl:  "stable url",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapToCert(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapToCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToCert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
