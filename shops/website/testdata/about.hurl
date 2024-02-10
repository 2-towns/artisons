# The about is displayed
GET {{host}}/about
HTTP 200
[Asserts]
xpath "//p[text()=\"About us\"]" exists
