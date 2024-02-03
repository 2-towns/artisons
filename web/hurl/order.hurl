# The blog article is displayed
GET {{host}}/account/orders/ORD1/detail.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//p[text()='ORD1']" exists
