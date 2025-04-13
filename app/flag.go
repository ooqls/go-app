package app

import "flag"

var registryPath string

var sqlFiles string

var RsaPrivKeyPath string
var RsaPubKeyPath string

var JwtPrivKeyPath string
var JwtPubKeyPath string


func init() {
	flag.StringVar(&registryPath, "registry", "/registry/registry.yaml", "Path to the registry path")
	flag.StringVar(&sqlFiles, "sql-files", "", "Comma separated list of files")
	flag.StringVar(&RsaPrivKeyPath, "rsa-private-key",  "/rsa/key.pem", "Path to an RSA private key")
	flag.StringVar(&RsaPubKeyPath, "rsa-public-key", "/rsa/pub.pem", "Path to the RSA public key")
	flag.StringVar(&JwtPrivKeyPath, "jwt-private-key", "/jwt/key.pem", "Path to a JWT private key")
	flag.StringVar(&JwtPubKeyPath, "jwt-public-key", "/jwt/pub.pem", "Path to a jwt public key")
}