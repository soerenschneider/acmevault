package client

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/acmevault/internal/config"
	"github.com/soerenschneider/acmevault/pkg/certstorage"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

type FSCertWriter struct {
	CertificateFile string
	PrivateKeyFile  string
	PemFile         string
	Uid             int
	Gid             int
}

func NewFsWriter(conf config.FsWriterConfig) (*FSCertWriter, error) {
	uid, err := getUidFromUsername(conf.Username)
	if err != nil {
		return nil, err
	}
	gid, err := getGidFromGroup(conf.Group)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("Resolved username %s to uid %d, group %s to %d", conf.Username, uid, conf.Group, gid)

	return &FSCertWriter{
		CertificateFile: conf.CertFile,
		PrivateKeyFile:  conf.PrivateKeyFile,
		PemFile:         conf.PemFile,
		Uid:             uid,
		Gid:             gid,
	}, nil
}

func getUidFromUsername(username string) (int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return -1, fmt.Errorf("could not fetch uid %s: %v", username, err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return -1, fmt.Errorf("can't cast uid %s", username)
	}

	return uid, nil
}

func getGidFromGroup(group string) (int, error) {
	g, err := user.LookupGroup(group)
	if err != nil {
		return -1, fmt.Errorf("could not fetch gid for group %s: %v", group, err)
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return -1, fmt.Errorf("can't cast gid %s", group)
	}

	return gid, nil
}

func (writer *FSCertWriter) ValidatePermissions() error {
	err := createIfNotExists(writer.CertificateFile, 0644, writer.Uid, writer.Gid)
	if err != nil {
		return fmt.Errorf("can not create cert file %s for user %d, %d: %v", writer.CertificateFile, writer.Uid, writer.Gid, err)
	}

	err = createIfNotExists(writer.PrivateKeyFile, 0600, writer.Uid, writer.Gid)
	if err != nil {
		return fmt.Errorf("can not create private file %s for user %d, %d: %v", writer.PrivateKeyFile, writer.Uid, writer.Gid, err)
	}

	return nil
}

func (writer *FSCertWriter) WriteBundle(cert *certstorage.AcmeCertificate) (bool, error) {
	if nil == cert {
		return false, errors.New("no certificate provided")
	}

	err := writer.ValidatePermissions()
	if err != nil {
		return false, fmt.Errorf("invalid permissions: %v", err)
	}

	runHook := false
	if len(writer.CertificateFile) > 0 {
		change, err := writeFile(writer.CertificateFile, cert.Certificate, 0640)
		if err != nil {
			return false, err
		}
		if change {
			runHook = true
		}
	}

	if len(writer.PrivateKeyFile) > 0 {
		change, err := writeFile(writer.PrivateKeyFile, cert.PrivateKey, 0600)
		if err != nil {
			return false, err
		}
		if change {
			runHook = true
		}
	}

	if len(writer.PemFile) > 0 {
		pem := []byte(cert.AsPem())
		change, err := writeFile(writer.PemFile, pem, 0600)
		if err != nil {
			return false, err
		}
		if change {
			runHook = true
		}
	}

	return runHook, nil
}

func writeFile(file string, content []byte, perms fs.FileMode) (bool, error) {
	identical, err := compareCerts(file, content)
	if err != nil || !identical {
		err = ioutil.WriteFile(file, content, perms)
		if err != nil {
			return false, fmt.Errorf("could not write pem to %s: %v", file, err)
		}
		return true, nil
	}
	return false, nil
}

func compareCerts(file string, payload []byte) (bool, error) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}

	return Md5Compare(fileContent, payload), nil
}

func Md5Compare(a, b []byte) bool {
	sumA := md5.Sum(a)
	sumB := md5.Sum(b)

	return bytes.Equal(sumA[:], sumB[:])
}

func createIfNotExists(path string, perm os.FileMode, uid, gid int) error {
	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("can not create file %s: %v", path, err)
	}

	err = os.Chown(path, uid, gid)
	if err != nil {
		fmt.Errorf("can not chown %s to %d %d: %v", path, uid, gid, err)
	}

	return nil
}
