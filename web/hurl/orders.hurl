# The orders list returns empty result
GET {{host}}/account/orders.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='order']" exists


