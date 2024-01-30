# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists


# Redirect to dashboard when the otp is correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email:hello@world.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Editing a locale without key shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
key: 
value: hop
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing a locale without value shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
key: hip
value: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#value-error" 
header "HX-Reswap" == "innerHTML show:#value-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing with correct data shows success
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
key: hip
value: hop
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 
