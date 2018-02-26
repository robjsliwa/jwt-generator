# jwt-generator
Simple utility for generating JWTs with claims listed in configuration file.  This utility can use existing private key or if the key is not provided it will generate new public/private keys.

# Configuration file
Configuration file uses yaml format and allows to specify token duration time, custom claims, and private/public key location.  Here is sample configuration.yaml:

```
expires: 129600

claims:
  - iss:jwt-generator
  - type:test

keys:
  public: ./jwt-public
  private: "./jwt-private"
```
