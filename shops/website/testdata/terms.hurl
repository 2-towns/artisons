# The terms page is displayed
GET {{host}}/terms-and-conditions
HTTP 200
[Asserts]
xpath "//p[text()=\"Terms and conditions\"]" exists
