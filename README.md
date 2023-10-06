# Banco Inter Go client

## Building command line tools

```sh
$ go build ./cmd/inter-token
$ go build ./cmd/inter-banking
```

## Authorize and get the user token

```
$ inter-token -scopes extrato.read -client-id <your client id> -client-secret <your client secret> --output-format info
token     a1200a94-b847-4cda-a510-cc0b9c7182d4
type      Bearer
expires   2022-02-02T02:22:22-03:00
scopes    extrato.read
```


## Use the banking tool

Show the command help using `inter-banking --help`

```
Usage: inter-banking [OPTION...] <COMMAND>

  -h, --help                 give this help list
  -c, --cert                 signed certificate file (default 'cert.crt')
  -k, --key                  certificate private key file (default 'cert.key')
  -t, --token                personal user token


balance                      get account balance

  -d, --date                 balance date in the format YYYY-MM-DD (defaults to
                             today)
  -f, --output-format        the output format used to show balance; can be
                             'short' (default) or 'full'

statement                    fetch account statements

  -s, --start-date           statements start date in the format YYYY-MM-DD
  -e, --end-date             statements end date in the format YYYY-MM-DD (defaults to
                             today)
```

### Fetch account balances

```
$ inter-banking --token a1200a94-b847-4cda-a510-cc0b9c7182d4 --date 2022-02-02 --output-format full
Balances at 2022-02-02

Available                  143293.57
Limit                           0.00
On hold                         0.00
Judicially blocked              0.00
Administratively blocked        0.00
```

### List account statements

```
$ inter-banking --token a1200a94-b847-4cda-a510-cc0b9c7182d4 --start-date 2022-02-02 --end-date 2023-02-12
Statements from 2022-02-02 to 2022-02-12

Date             Value   Operation  Type           Title                   Description
2022-02-02    22373.32   credit     transferencia  Transferência recebida  TED RECEBIDA - 001 BANCO 001 S.A.            
2022-02-05    22300.00   debit      pix            Pix enviado             PIX ENVIADO - Cp :123456  
2022-02-09     1000.00   debit      pix            Pix enviado             PIX ENVIADO - Cp :789012
```

