# Redirect to dashboard when the otp is correct
POST {{host}}/login
HX-Current-Url: http://localhost:8080/sso
HX-Request: true
[FormParams]
email: login-admin@artisons.me
otp: 1
otp: 2
otp: 3 
otp: 4 
otp: 5
otp: 6 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index"
cookie "wsid" exists
cookie "wsid[HttpOnly]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[Max-Age]" exists
cookie "wsid[SameSite]" == "Strict"
cookie "wsid[Path]" == "/"