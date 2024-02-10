# The orders list returns empty result
GET {{host}}/account/orders
[Cookies]
wsid: 333333
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='order']" exists


