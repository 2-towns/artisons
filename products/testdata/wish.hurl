# The wish list returns empty result
GET {{host}}/account/wish
[Cookies]
wsid: 333333
HTTP 200
[Asserts]
body contains "No results found."

# Add item to the wish list
POST {{host}}/account/wish/PDT1/add
HX-Request: true
[Cookies]
wsid: 333333
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='wishlist-card']" exists

# Delete item to the wish list
POST {{host}}/account/wish/PDT1/delete
HX-Request: true
[Cookies]
wsid: 333333
HTTP 200
[Asserts]
body contains "No results found."

