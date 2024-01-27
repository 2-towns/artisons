# Redirect to auth page when the user is not connected
GET {{host}}/account/index.html
HTTP 302
[Asserts]
header "Location" == "/otp.html"

# Returns error when trying to post a non htmx request
POST {{host}}/otp
[FormParams]
email: arnaud@yandex.com
HTTP 400
[Asserts]
body contains "Your are not authorized to access to this page. This will be reported."

# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: arnaud@yandex.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists
[Captures]
glue: xpath "normalize-space(//input[@id='glue']/@value)"

# Displays an error message when requesting new OTP too quickly
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: arnaud@yandex.com
HTTP 200
[Asserts]
body contains "You need to wait before asking another OTP."

# Displays an error when otp is not correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
glue: {{glue}}
otp: 1
otp: 2 
otp: 3 
otp: 4 
otp: 5
otp: 6 
HTTP 200
[Asserts]
body contains "The OTP does not match."

# Redirect to account when the otp is correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
glue: {{glue}}
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index.html"

# Redirect when the user is connected and try to access to auth
GET {{host}}/otp.html
HTTP 302
[Asserts]
header "Location" == "/account/index.html"