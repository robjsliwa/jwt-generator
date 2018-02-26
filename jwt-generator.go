package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	yaml "gopkg.in/yaml.v2"
)

// Config - configuration file with JWT specs
type Config struct {
	Expires int64
	Claims  []string
	Keys    struct {
		Public  string
		Private string
	}
}

func generatePublicPrivateKeys() (string, string, error) {
	// RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// get private key
	privateKeyDerEncoded := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDerEncoded,
	}
	privateKeyPem := string(pem.EncodeToMemory(&privateKeyBlock))

	// get public key
	publicKey := privateKey.PublicKey
	publicKeyDerEncoded, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDerEncoded,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))

	return privateKeyPem, publicKeyPem, nil
}

func generateJWT(conf *Config, privateKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	tokenDuration := time.Duration(conf.Expires)

	claims["exp"] = time.Now().Add(time.Minute * tokenDuration).Unix()
	claims["iat"] = time.Now().Unix()

	for _, claim := range conf.Claims {
		claimKV := strings.Split(claim, ":")
		if len(claimKV) != 2 {
			return "", fmt.Errorf("invalid claim: %v", claim)
		}
		claimName := claimKV[0]
		claimValue := claimKV[1]
		claims[claimName] = claimValue
	}

	privateKeyBytes, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		fmt.Printf("Error while parsing private key from PEM: %v\n", err)
		fmt.Printf("Received key: %v\n", privateKey)
		return "", err
	}

	tokenString, err := token.SignedString(privateKeyBytes)
	if err != nil {
		fmt.Printf("Error while signing JWT: %v\n", err)
		return "", err
	}

	return tokenString, nil
}

func main() {
	configFile := flag.String("config", "config.yaml", "Configuration file with JWT specs")
	flag.Parse()

	var conf Config
	r, err := os.Open(*configFile)
	if err != nil {
		log.Fatalf("Config file error: %v\n", err)
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("Config open error: %v\n", err)
	}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		log.Fatalf("Could not parse config: %v\n", err)
	}
	fmt.Printf("%+v\n", conf)

	// check if private key exists, if it does not generate new keys
	var privateKey string
	var publicKey string
	if _, err = os.Stat(conf.Keys.Private); os.IsNotExist(err) {
		fmt.Println("No keys found, will generate new keys")
		var errKey error
		privateKey, publicKey, errKey = generatePublicPrivateKeys()
		if errKey != nil {
			log.Fatalf("Generate keys error: %v\n", errKey)
		}
		fileErr := ioutil.WriteFile(conf.Keys.Private, []byte(privateKey), 0644)
		if fileErr != nil {
			log.Fatalf("Error saving private key: %v\n", fileErr)
		}
		fileErr = ioutil.WriteFile(conf.Keys.Public, []byte(publicKey), 0644)
		if fileErr != nil {
			log.Fatalf("Error saving public key: %v\n", fileErr)
		}
	} else {
		fmt.Printf("Loading private key: %v...\n", conf.Keys.Private)
		data, fileErr := ioutil.ReadFile(conf.Keys.Private)
		if fileErr != nil {
			log.Fatalf("Error reading private file: %v\n", fileErr)
		}
		privateKey = string(data)
	}

	token, err := generateJWT(&conf, privateKey)
	if err != nil {
		log.Fatalf("Generate JWT failed: %v\n", err)
	}

	fmt.Printf("Token: %v\n", token)
}
