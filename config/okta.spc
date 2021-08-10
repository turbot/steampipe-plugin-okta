# connection "okta" {
  plugin = "okta"

  # 1. With API TOKEN(https://developer.okta.com/docs/guides/create-an-api-token/create-the-token/)
  # domain = "https://<your_okta_domain>.okta.com"
  # token  = "this_not_real_token"


  # 2. With Private Key(https://github.com/okta/okta-sdk-golang#oauth-20)
  # domain      = "https://<your_okta_domain>.okta.com"
  # client_id   = "your_client_id"
  # private_key = "----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAmyX8wdrHK1ycOMeXNg3NOMQvebnfQp+3L5OaaiX16/+tLbwb\nJTZDYh0EXLySMVsduRxC/1PQdPuI6x50TdkoB3C4JMuU968uJqkFp7fXXy5SMAej\nHAyF67cY51dx15ztvakRNJPhhI5WaC20RfR/eow0IH5lGI3czcvTCChGau5qLue3\nHqNDYFY+U3xhOlavSDdtmuxpIFsDycn/OjYjsV4lzyRrOArqtVV/kXHKx04T6A1x\nSc99999999999999999999999999999999999999999999999999EGekHlUAIUpw\n-----END RSA PRIVATE KEY-----"
}
