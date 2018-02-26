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

# How to use
./jwt-generator --configuration sample-jwt.yaml

This produces following output:

```
No keys found, will generate new keys
Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Mjc0Mzk2NDYsImlhdCI6MTUxOTY2MzY0NiwiaXNzIjoiand0LWdlbmVyYXRvciIsInR5cGUiOiJ0ZXN0In0.JVplOSdKeeyypPPAxYDLmpUEVTjunwuVPRonlqdlcHYWiS1ssw7xRlUwwvsB1sRAnf-aXBeBoOvPosswoHDSm7AlrWjMAIONw0_PgI1TPluk_TEr5_syG1uVMNDL7QChKAVO0tpW1eJoa1KBhb0WU4we6gw_FpJBpdn2piXdDTJf35U_AWqSbmXzcy1eZy0-VAcROnER4QS7ujjCQZV5LjQD1p0zRbMzGoANR2RPU6C2VFvWjABAIXhcwuaxZX65YNOtABgcuKXjMZZOGgzcHiNbTXNyfBH4FAEdqfKnyulr_DkJKVjEr8c94RHL7kBEkyycSh62mvhDJDts8n-xmw
```
