# Redirect to auth page when the user is not connected
GET {{host}}/account/index.html
HTTP 302
[Asserts]
header "Location" == "/otp.html"

# Returns error when trying to post a non htmx request
POST {{host}}/otp.html
[FormParams]
email: hello@artisons.me
HTTP 400
[Asserts]
body contains "Your are not authorized to access to this page. This will be reported."

# Displays an error when otp is not correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: arnaud@artisons.me
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
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: arnaud@artisons.me
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index.html"
cookie "wsid" exists
cookie "wsid[HttpOnly]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[SameSite]" == "Strict"
cookie "wsid[Path]" == "/"

# Redirect when the admin is connected and try to access to auth
GET {{host}}/otp.html
HX-Current-URL: http://localhost:8080/otp.html
HTTP 302
[Asserts]
header "Location" == "/account/index.html"