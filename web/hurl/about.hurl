# The about is displayed
GET {{host}}/about.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//p[text()=\"About us\"]" exists
