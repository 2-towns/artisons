# The wish list returns empty result
GET {{host}}/account/wish.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."

# Add item to the wish list
POST {{host}}/account/wish/PDT1/add.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='wishlist-card']" exists

# Delete item to the wish list
POST {{host}}/account/wish/PDT1/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."

