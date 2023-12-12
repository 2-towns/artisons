# Redirect to auth page when the user is not connected
GET {{host}}/admin/dashboard.html
HTTP 302
[Asserts]
header "Location" == "/auth"

# Returns error when trying to post a non htmx request
POST {{host}}/otp
[FormParams]
email: hello@world.com
HTTP 400
[Asserts]
body contains "The request is invalid. This will be reported."

# Displays an inline error when the email is not valid
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: arnaud@yandex.com
HTTP 200
[Asserts]
body contains "<div class=\"form-error\">The email is invalid.</div>"

# Display otp form when the correct email is used
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/auth/login.html']" exists
[Captures]
glue: xpath "normalize-space(//input[@id='glue']/@value)"

# Displays an error message when requesting new OTP too quickly
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
body contains "You need to wait before asking another OTP."

# Displays an error when otp is not correct
POST {{host}}/auth/login.html
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

# Redirect to dashboard when the otp is correct
POST {{host}}/auth/login.html
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
header "HX-Redirect" == "/admin"

# Redirect when the admin is connected and try to access to auth
GET {{host}}/auth
HTTP 302
[Asserts]
header "Location" == "/admin"