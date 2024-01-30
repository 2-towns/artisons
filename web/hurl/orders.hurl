# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: arnaud@yandex.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists

# Redirect to account when the otp is correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: arnaud@yandex.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index.html"

# The orders list returns empty result
GET {{host}}/account/orders.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='order']" exists


