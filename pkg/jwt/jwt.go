package jwt

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type customClaims struct {
	Payload []byte `json:"payload"`
	jwt.RegisteredClaims
}

func GenTokenAlgHS256(issuer string, secret string, payload any, expireTime time.Time) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(payload); err != nil {
		return "", errors.WithMessage(err, "enc.Encode fail")
	}

	claims := customClaims{
		buf.Bytes(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.WithMessage(err, "jwt.SignedString")
	}

	return ss, nil
}

func ParseTokenAlgHS256(tokenString string, secret string, payload any) error {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return errors.WithMessage(err, "jwt.Parse")
	}

	c, ok := token.Claims.(*customClaims)
	if !ok {
		return errors.New("invalid token")
	}

	buf := bytes.NewBuffer(c.Payload)

	dec := gob.NewDecoder(buf)

	if err = dec.Decode(payload); err != nil {
		return errors.WithMessage(err, "gob.Decode")
	}

	return nil
}

func DecodeSegment(tokenString string, payload any) error {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return errors.New("invalid token")
	}

	payloadPart := parts[1]
	payloadBytes, err := jwt.NewParser().DecodeSegment(payloadPart)
	if err != nil {
		return errors.WithMessage(err, "jwt.NewParser().DecodeSegment")
	}

	claims := &customClaims{}

	if err = json.Unmarshal(payloadBytes, claims); err != nil {
		return errors.WithMessage(err, "json.Unmarshal")
	}

	buf := bytes.NewBuffer(claims.Payload)

	dec := gob.NewDecoder(buf)

	if err = dec.Decode(payload); err != nil {
		return errors.WithMessage(err, "gob.Decode")
	}

	return nil
}

// GenTokenAlgRS256
// privateKey 產生方式如下：
// signKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes.NewBufferString(privateKey).Bytes())
// if err != nil {
// return "", errors.Wrapf(err, "jwt.ParseRSAPrivateKeyFromPEM")
// }
func GenTokenAlgRS256(issuer string, privateKey *rsa.PrivateKey, payload any, expireTime time.Time) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(payload); err != nil {
		return "", errors.WithMessage(err, "enc.Encode fail")
	}

	claims := customClaims{
		buf.Bytes(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.WithMessage(err, "jwt.SignedString")
	}

	return ss, nil
}

func ParseTokenAlgRS256(tokenString string, publicKey *rsa.PublicKey, payload any) error {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return errors.WithMessage(err, "jwt.Parse")
	}

	c, ok := token.Claims.(*customClaims)
	if !ok {
		return errors.New("invalid token")
	}

	buf := bytes.NewBuffer(c.Payload)

	dec := gob.NewDecoder(buf)

	if err = dec.Decode(payload); err != nil {
		return errors.WithMessage(err, "gob.Decode")
	}

	return nil
}
