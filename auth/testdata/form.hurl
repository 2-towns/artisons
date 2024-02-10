# The otp form shows the email input 
GET {{host}}/otp
HTTP 200
[Asserts]
xpath "//input[@name='email']" exists 

# The sso form shows the email input 
GET {{host}}/sso
HTTP 200
[Asserts]
xpath "//input[@name='email']" exists 

# The user is redirected when signed in
GET {{host}}/otp
[Cookies]
wsid: 333333
HTTP 302
[Asserts]
header "Location" == "/account/index"

# The user is redirected when signed in
GET {{host}}/sso
[Cookies]
wsid: 333333
HTTP 302
[Asserts]
header "Location" == "/admin/index"

