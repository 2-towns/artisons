# Redirect to auth page when the user is not connected
GET {{host}}/admin/dashboard.html
HTTP 302
[Asserts]
header "Location" == "/sso.html"

# Returns error when trying to post a non htmx request
POST {{host}}/otp.html
[FormParams]
email: hello@world.com
HTTP 400
[Asserts]
body contains "Your are not authorized to access to this page. This will be reported."

# Displays an inline error when the email is not valid
POST {{host}}/otp.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
[FormParams]
email: arnaud@yandex.com
HTTP 200
[Asserts]
body contains "<div class=\"form-error\">The data is invalid.</div>"

# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists

# Displays an error message when requesting new OTP too quickly
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
body contains "You need to wait before asking another OTP."

# Displays an error when otp is not correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: hello@world.com
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
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: hello@world.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Redirect when the admin is connected and try to access to auth
GET {{host}}/sso.html
HX-Current-URL: http://localhost:8080/sso.html
HTTP 302
[Asserts]
header "Location" == "/admin/index.html"