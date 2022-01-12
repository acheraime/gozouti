package tls

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

type ParsedCertificate struct {
	Certificate   *x509.Certificate
	Intermediates []*x509.Certificate
}

func ParseCertificate(c []byte) (*ParsedCertificate, error) {
	// certs := x509.NewCertPool()
	// ok := certs.AppendCertsFromPEM(c)
	// if !ok {
	// 	return fmt.Errorf("unable to parse certificate")
	// }
	var parsedCert ParsedCertificate
	certs, err := DecodePemChain(c)
	if err != nil {
		return nil, err
	}
	for _, cert := range certs {
		if !cert.IsCA {
			parsedCert.Certificate = cert
			continue
		}
		parsedCert.Intermediates = append(parsedCert.Intermediates, cert)

	}
	//log.Println(cert.DNSNames, cert.Issuer, cert.NotBefore, cert.NotAfter, cert.Subject, cert.SerialNumber)
	// intermediates := x509.NewCertPool()
	// intermediates.AddCert(certs[2])
	// intermediates.AddCert(certs[3])

	// verify certificate
	// verfiyOptions := x509.VerifyOptions{
	// 	KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	// 	Roots:         x509.NewCertPool(),
	// 	Intermediates: intermediates,
	// }
	// _, err = certs[0].Verify(verfiyOptions)
	// if err != nil {
	// 	return err
	// }

	return &parsedCert, nil
}

func DecodePemChain(chain []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	var derBlock *pem.Block
	for {
		derBlock, chain = pem.Decode(chain)
		if derBlock == nil {
			break
		}

		if derBlock.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(derBlock.Bytes)
			if err != nil {
				break
			}
			certs = append(certs, cert)
		}
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("unable to parse pem chain")
	}

	return certs, nil
}

func (c ParsedCertificate) IsValid() bool {
	return c.Certificate.NotAfter.After(time.Now())
}

func (c ParsedCertificate) ValidityPeriod() string {
	validFrom := c.Certificate.NotBefore.Format(time.RFC822)
	validTo := c.Certificate.NotAfter.Format(time.RFC822)

	return validFrom + " - " + validTo
}
