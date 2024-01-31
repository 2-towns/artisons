# The terms page is displayed
GET {{host}}/terms-and-conditions.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//p[text()=\"Terms and conditions\"]" exists
