# Redirect to the login form  when the admin is not signed in 
GET {{host}}/admin/index
HTTP 302
[Asserts]
header "Location" == "/sso"

# Redirect to the login form  when the user is not signed in 
GET {{host}}/account/index
HTTP 302
[Asserts]
header "Location" == "/otp"

# Remove session id when it does not exist 
GET {{host}}/account/index
[Cookies]
wsid: 09090909090 
HTTP 302
[Asserts]
header "Location" == "/otp"
cookie "wsid" exists
cookie "wsid[Max-Age]" == 0

# Refresh the user session id
GET {{host}}/account/index
[Cookies]
wsid: 333333 
HTTP 200
[Asserts]
cookie "wsid" exists
cookie "wsid[Max-Age]" > 0