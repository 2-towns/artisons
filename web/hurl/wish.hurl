# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: arnaud{{time}}@yandex.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists

# Redirect to account when the otp is correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: arnaud{{time}}@yandex.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index.html"

# The wish list returns empty result
GET {{host}}/wish.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."

# Add item to the wish list
POST {{host}}/wish/PDT1/add.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//div[@class='wishlist-card']" exists

# Delete item to the wish list
POST {{host}}/wish/PDT1/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."

