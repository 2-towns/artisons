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

# The blog article is displayed
GET {{host}}/account/orders/ORD1/detail.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//p[text()='ORD1']" exists
