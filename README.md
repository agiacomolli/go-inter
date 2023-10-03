# Banco Inter Go client

## Building command line tools

```sh
$ go build ./cmd/inter-token
$ go build ./cmd/inter-banking
```

## Authorize and get the user token

```sh
$ ./inter-token -scopes extrato.read -client-id <your client id> -client-secret <your client secret> --output-format info
token     a1200a94-b847-4cda-a510-cc0b9c7182d4
type      Bearer
expires   2022-02-02T02:22:22-03:00
scopes    extrato.read
```

## List account statements

```sh
./inter-banking -token a1200a94-b847-4cda-a510-cc0b9c7182d4 -start 2022-02-02 -end 2023-02-12
Statements from 2022-02-02 to 2022-02-12

Date             Value   Operation  Type           Title                   Description
2022-02-02    22373.32   credit     transferencia  Transferência recebida  TED RECEBIDA - 001 BANCO 001 S.A.            
2022-02-05    22300.00   debit      pix            Pix enviado             PIX ENVIADO - Cp :123456  
2022-02-09     1000.00   debit      pix            Pix enviado             PIX ENVIADO - Cp :789012
```
