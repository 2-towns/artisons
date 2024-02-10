# Displays an error when email is empty
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#email-error" 
header "HX-Reswap" == "innerHTML show:#email-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Displays an error when email is invalid
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: invalid
HTTP 200
[Asserts]
header "HX-Retarget" == "#email-error" 
header "HX-Reswap" == "innerHTML show:#email-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Show user login form
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[FormParams]
email: user@artisons.me
HTTP 200
[Asserts]
xpath "//a[@href='/otp']" exists

# Show admin login form
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/sso
HX-Request: true
[FormParams]
email: admin@artisons.me
HTTP 200
[Asserts]
xpath "//a[@href='/sso']" exists

# The user is redirected when signed in
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/otp
HX-Request: true
[Cookies]
wsid: 333333
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index"