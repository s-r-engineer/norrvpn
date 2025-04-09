package libraryEncryption

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	libraryNetwork "github.com/s-r-engineer/library/network"
)

func GetDHSecretFromConnection(conn libraryNetwork.GenericConnection, p *big.Int, g *big.Int) (string, error) {
	priv, err := rand.Int(rand.Reader, p)
	if err != nil {
		return "", err
	}
	pub := new(big.Int).Exp(g, priv, p)
	otherSidePub := make([]byte, p.BitLen()/8+1)
	_, err = conn.Write(pub.Bytes())
	if err != nil {
		return "", err
	}
	n, err := conn.Read(otherSidePub)
	if err != nil {
		return "", err
	}
	otherSide := new(big.Int).SetBytes(otherSidePub[:n])
	sharedSecret := new(big.Int).Exp(otherSide, priv, p)
	symmetricKey := sha256.Sum256(sharedSecret.Bytes())
	return fmt.Sprintf("%x", symmetricKey), nil
}
