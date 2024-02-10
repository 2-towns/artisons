# Logout admin redirect to login form
POST {{host}}/logout
HX-Request: true
HX-Current-Url: /admin/index
[Cookies]
wsid: 987654321
HTTP 200
[Asserts]
header "HX-Redirect" == "/sso"
cookie "wsid" exists
cookie "wsid[Max-Age]" == 0

# The admin is not authorized to see the dashboard
GET {{host}}/admin/index
[Cookies]
wsid: 987654321
HTTP 302
[Asserts]
header "Location" == "/sso"

# Logout with bad session show not authorized message
POST {{host}}/logout
HX-Request: true
HX-Current-Url: /admin/index
[Cookies]
wsid: idontexist
HTTP 200
[Asserts]
xpath "normalize-space(//p[@class='alert-message']/text())" == "Your are not authorized to access to this page. This will be reported." 
cookie "wsid" exists
cookie "wsid[Max-Age]" == 0

# Access without the cookie redirect to the login page 
POST {{host}}/logout
HX-Request: true
HX-Current-Url: /admin/index
HTTP 200
[Asserts]
header "HX-Redirect" == "/sso"

# Logout user redirect to home form
POST {{host}}/logout
HX-Request: true
[Cookies]
wsid: 222222
HTTP 200
[Asserts]
header "HX-Redirect" == "/"
cookie "wsid" exists
cookie "wsid[Max-Age]" == 0

# The admin is not authorized to see the account page
GET {{host}}/account/index
[Cookies]
wsid: 222222
HTTP 302
[Asserts]
header "Location" == "/otp"

# Access without the cookie redirect to the home page 
POST {{host}}/logout
HX-Request: true
HX-Current-Url: /
HTTP 200
[Asserts]
header "HX-Redirect" == "/"
