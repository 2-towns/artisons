# The french addresse api gives results
GET {{host}}/addresses?q=8+bd+du+port
HTTP 200
[Asserts]
xpath "//li[@class='address-item']" exists
