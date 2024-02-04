# Redirect to auth page when the user is not connected
GET {{host}}/admin/dashboard.html
HTTP 302
[Asserts]
header "Location" == "/sso.html"

# Returns error when trying to post a non htmx request
POST {{host}}/otp.html
[FormParams]
email: hello@artisons.me
HTTP 400
[Asserts]
body contains "Your are not authorized to access to this page. This will be reported."

# Displays an inline error when the email is not valid
POST {{host}}/otp.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
[FormParams]
email: arnaud@artisons.me
HTTP 200
[Asserts]
body contains "<div class=\"form-error\">The data is invalid.</div>"

# Displays an error when otp is not correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: hello@artisons.me
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
email: hello@artisons.me
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"
cookie "wsid" exists
cookie "wsid[HttpOnly]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[SameSite]" == "Strict"
cookie "wsid[Path]" == "/"

# Redirect when the admin is connected and try to access to auth
GET {{host}}/sso.html
HX-Current-URL: http://localhost:8080/sso.html
HTTP 302
[Asserts]
header "Location" == "/admin/index.html"