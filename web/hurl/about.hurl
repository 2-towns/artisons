# The about is displayed
GET {{host}}/about.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//p[text()=\"About us\"]" exists
