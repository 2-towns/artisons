# The blog article is displayed
GET {{host}}/account/orders/ORD1/detail
[Cookies]
wsid: 333333
HTTP 200
[Asserts]
xpath "//p[text()='ORD1']" exists
