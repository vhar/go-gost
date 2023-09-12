package signature

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ftomza/gogost/gost28147"
	"github.com/ftomza/gogost/gost34112012256"
	"github.com/ftomza/gogost/gost34112012512"
	"github.com/ftomza/gogost/gost341194"
	"go.mozilla.org/pkcs7"
)

type Sign struct {
	p7   *pkcs7.PKCS7
	data []byte
}

type Signer struct {
	CommonName          string
	CountryName         string
	StateOrProvinceName string
	LocalityName        string
	Surname             string
	GivenName           string
	Title               string
	EmailAddress        string
}

type Certificate struct {
	IssuerName          string
	NotBefore           string
	NotAfter            string
	EncriptionAlgorithm string
	DigestAlgorithm     string
}

// Структура отчета о проверке ЭЦП
type Report struct {
	Signer      Signer
	Certificate Certificate
	SigningTime string
	Validity    bool
}

var (
	OIDEmailAddress           = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
	OIDAttributeMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
	OIDAttributeSigningTime   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}
	OIDgostR341194            = asn1.ObjectIdentifier{1, 2, 643, 2, 2, 9}
	OIDgostR34102001          = asn1.ObjectIdentifier{1, 2, 643, 2, 2, 19}
	OIDgostR34102012256       = asn1.ObjectIdentifier{1, 2, 643, 7, 1, 1, 1, 1}
	OIDgostR34102012512       = asn1.ObjectIdentifier{1, 2, 643, 7, 1, 1, 1, 2}
	OIDgostR34112012256       = asn1.ObjectIdentifier{1, 2, 643, 7, 1, 1, 2, 2}
	OIDgostR34112012512       = asn1.ObjectIdentifier{1, 2, 643, 7, 1, 1, 2, 3}
	OIDCommonName             = asn1.ObjectIdentifier{2, 5, 4, 3}
	OIDSurname                = asn1.ObjectIdentifier{2, 5, 4, 4}
	OIDCountryName            = asn1.ObjectIdentifier{2, 5, 4, 6}
	OIDLocalityName           = asn1.ObjectIdentifier{2, 5, 4, 7}
	OIDStateOrProvinceName    = asn1.ObjectIdentifier{2, 5, 4, 8}
	OIDStreetAddress          = asn1.ObjectIdentifier{2, 5, 4, 9}
	OIDOrganizationName       = asn1.ObjectIdentifier{2, 5, 4, 10}
	OIDOrganizationUnitName   = asn1.ObjectIdentifier{2, 5, 4, 11}
	OIDTitle                  = asn1.ObjectIdentifier{2, 5, 4, 12}
	OIDGivenName              = asn1.ObjectIdentifier{2, 5, 4, 42}
)

// Декодируем содержимое файла подписи
// Если возникла ошибка при декодировании мз Base64, считаем что подпись в формате PKCS7
func NewSign(data []byte, sign []byte) (*Sign, error) {
	signature, err := base64.StdEncoding.DecodeString(string(sign))
	if err != nil {
		signature = sign
	}

	p7, err := pkcs7.Parse(signature)
	if err != nil {
		return nil, fmt.Errorf("Ошибка декодирования файла подписи")
	}

	return &Sign{
		p7:   p7,
		data: data,
	}, nil
}

// Определяем алгоритм хэширования при помощи которого создавалась ЭЦП
func (s *Sign) GetHash() ([]byte, error) {
	digestAlgorithm := s.p7.Signers[0].DigestAlgorithm.Algorithm

	if asn1.ObjectIdentifier.Equal(digestAlgorithm, OIDgostR341194) {
		hasher := gost341194.New(&gost28147.SboxIdGostR341194TestParamSet)
		_, err := hasher.Write(s.data)
		if err != nil {
			return nil, err
		}

		return hasher.Sum(nil), nil

	} else if asn1.ObjectIdentifier.Equal(digestAlgorithm, OIDgostR34112012256) {
		hasher := gost34112012256.New()
		_, err := hasher.Write(s.data)
		if err != nil {
			return nil, err
		}

		return hasher.Sum(nil), nil

	} else if asn1.ObjectIdentifier.Equal(digestAlgorithm, OIDgostR34112012512) {
		hasher := gost34112012512.New()
		_, err := hasher.Write(s.data)
		if err != nil {
			return nil, err
		}

		return hasher.Sum(nil), nil
	}

	return nil, fmt.Errorf("Невозможно определить алгоритм хэширования")
}

// Проверка вычисленного хэша документа и значения из файла подписи
// Если данные строки совпадают, значит файл подписи сделан для данного документа и документ после этого не изменялся
func (s *Sign) Verify() (bool, error) {
	dgst, err := s.GetHash()
	if err != nil {
		return false, fmt.Errorf("Ошибка записи данных в функцию хэширования. %s", err.Error())
	}

	var messageDigest []byte

	if err = s.p7.UnmarshalSignedAttribute(OIDAttributeMessageDigest, &messageDigest); err != nil {
		return false, fmt.Errorf("Ошибка десериализации OID хэша. %s", err.Error())
	}

	return hex.EncodeToString(dgst) == hex.EncodeToString(messageDigest), nil
}

// Выборка данных о сертификате ЭЦП
func (s *Sign) GetCertificate() Certificate {
	cert := Certificate{
		IssuerName: s.p7.GetOnlySigner().Issuer.CommonName,
		NotBefore:  s.p7.GetOnlySigner().NotBefore.Format("02.01.2006 15:04"),
		NotAfter:   s.p7.GetOnlySigner().NotAfter.Format("02.01.2006 15:04"),
	}

	encryptionAlgorithm := s.p7.Signers[0].DigestEncryptionAlgorithm.Algorithm

	if asn1.ObjectIdentifier.Equal(encryptionAlgorithm, OIDgostR34102001) {
		cert.EncriptionAlgorithm = "ГОСТ Р 34.11/34.10-2001"
	} else if asn1.ObjectIdentifier.Equal(encryptionAlgorithm, OIDgostR34102012256) {
		cert.EncriptionAlgorithm = "ГОСТ Р 34.11/34.10-2012 (256 бит)"
	} else if asn1.ObjectIdentifier.Equal(encryptionAlgorithm, OIDgostR34102012512) {
		cert.EncriptionAlgorithm = "ГОСТ Р 34.11/34.10-2012 (512 бит)"
	}

	digestAlgorithm := s.p7.Signers[0].DigestAlgorithm.Algorithm

	if asn1.ObjectIdentifier.Equal(digestAlgorithm, OIDgostR341194) {
		cert.DigestAlgorithm = "ГОСТ Р 34.11-94"
	} else if asn1.ObjectIdentifier.Equal(digestAlgorithm, OIDgostR34112012256) {
		cert.DigestAlgorithm = "ГОСТ Р 34.11-2012 (256 бит)"
	} else if asn1.ObjectIdentifier.Equal(digestAlgorithm, OIDgostR34112012512) {
		cert.DigestAlgorithm = "ГОСТ Р 34.11-2012 (512 бит)"
	}

	return cert
}

// Выборка сведений о подписанте
func (s *Sign) GetSigner() Signer {
	signer := Signer{}

	names := s.p7.GetOnlySigner().Subject.Names
	for _, attr := range names {
		if asn1.ObjectIdentifier.Equal(attr.Type, OIDCommonName) {
			signer.CommonName = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDSurname) {
			signer.Surname = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDCountryName) {
			signer.CountryName = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDLocalityName) {
			signer.LocalityName = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDStateOrProvinceName) {
			signer.StateOrProvinceName = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDTitle) {
			signer.Title = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDGivenName) {
			signer.GivenName = attr.Value.(string)
		} else if asn1.ObjectIdentifier.Equal(attr.Type, OIDEmailAddress) {
			signer.EmailAddress = attr.Value.(string)
		}
	}

	return signer
}

// Получение сведений о времени подписания документа
// Значение данного поля несет информативный характер, так как время подпси берется из системы подписанта
func (s *Sign) GetSigningTume() (signingTime time.Time, err error) {
	signingTime = time.Now().UTC()

	err = s.p7.UnmarshalSignedAttribute(OIDAttributeSigningTime, &signingTime)
	if err != nil {
		return signingTime, fmt.Errorf("Error unmarshal signed time OID. %s", err.Error())
	}
	return signingTime, nil
}

func (s *Sign) GetReport() Report {
	report := Report{
		Signer:      s.GetSigner(),
		Certificate: s.GetCertificate(),
	}

	report.Validity, _ = s.Verify()

	signingTime, _ := s.GetSigningTume()
	report.SigningTime = signingTime.Format("02.01.2006 15:04:05 MST")

	return report
}
