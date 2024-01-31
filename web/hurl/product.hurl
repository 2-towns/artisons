# The offline product returns 404
GET {{host}}/product/offline.html
User-Agent: {{ua}}
HTTP 404
[Asserts]
xpath "//h2[text()=\"Oops..You just found an error page.\"]" exists
