# Displays an inline error when the email is not valid
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: 
HTTP 200
[Asserts]
body contains "<div class=\"form-error\">The data is invalid.</div>"

# Displays an error when otp is not correct
POST {{host}}/login
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: otp@artisons.me
otp: 1
otp: 1 
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
body contains "The OTP does not match."

# Displays an error when email does not have otp
POST {{host}}/login
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: idontexist@artisons.me
otp: 1
otp: 1 
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
xpath "normalize-space(//p[@class='alert-message']/text())" == "Your are not authorized to access to this page. This will be reported." 

# Redirect to account when the otp is correct
POST {{host}}/login
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: otp@artisons.me
otp: 1
otp: 2
otp: 3 
otp: 4 
otp: 5
otp: 6 
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index"
cookie "wsid" exists
cookie "wsid[HttpOnly]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[SameSite]" == "Strict"
cookie "wsid[Path]" == "/"
