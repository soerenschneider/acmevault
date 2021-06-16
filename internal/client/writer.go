package client

import (
	"acmevault/pkg/certstorage"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

type FSCertWriter struct {
	CertificatePath string
	PrivateKeyPath  string
	Uid             int
	Gid             int
}

func NewFsWriter(certPath, privateKeyPath string, username, group string) (*FSCertWriter, error) {
	uid, err := getUidFromUsername(username)
	if err != nil {
		return nil, err
	}
	gid, err := getGidFromGroup(group)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("Resolved username %s to uid %d, group %s to %d", username, uid, group, gid)

	return &FSCertWriter{
		CertificatePath: certPath,
		PrivateKeyPath:  privateKeyPath,
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
	err := createIfNotExists(writer.CertificatePath, 0644, writer.Uid, writer.Gid)
	if err != nil {
		return fmt.Errorf("can not create cert file %s for user %d, %d: %v", writer.CertificatePath, writer.Uid, writer.Gid, err)
	}

	err = createIfNotExists(writer.PrivateKeyPath, 0600, writer.Uid, writer.Gid)
	if err != nil {
		return fmt.Errorf("can not create private file %s for user %d, %d: %v", writer.PrivateKeyPath, writer.Uid, writer.Gid, err)
	}

	return nil
}

func (writer *FSCertWriter) WriteBundle(cert *certstorage.AcmeCertificate) (bool, error) {
	if nil == cert {
		return false, errors.New("Empty certificate provided")
	}

	writer.ValidatePermissions()

	runHook := false
	identical, err := compareCerts(writer.CertificatePath, cert.Certificate)
	if err != nil || !identical {
		err := ioutil.WriteFile(writer.CertificatePath, cert.Certificate, 0644)
		if err != nil {
			log.Fatal().Msgf("could not writeCertificate cert to %s: %v", writer.CertificatePath, err)
		}

		runHook = true
	}

	identical, err = compareCerts(writer.PrivateKeyPath, cert.PrivateKey)
	if err != nil || !identical {
		err = ioutil.WriteFile(writer.PrivateKeyPath, cert.PrivateKey, 0600)
		if err != nil {
			log.Fatal().Msgf("could not private key to %s: %v", writer.PrivateKeyPath, err)
		}

		runHook = true
	}

	return runHook, nil
}

func compareCerts(path string, payload []byte) (bool, error) {
	fileContent, err := ioutil.ReadFile(path)
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
