# Display otp form when the correct email is used
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/auth/login.html']" exists
[Captures]
glue: xpath "normalize-space(//input[@id='glue']/@value)"

# Redirect to dashboard when the otp is correct
POST {{host}}/auth/login.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
glue: {{glue}}
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Editing a locale without locale shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
locale:
key: hip
value: hop
HTTP 200
[Asserts]
header "HX-Retarget" == "#locale-error" 
header "HX-Reswap" == "innerHTML show:#locale-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The locale is invalid." 

# Editing a locale with a bad locale shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
locale: !!!
key: hip
value: hop
HTTP 200
[Asserts]
header "HX-Retarget" == "#locale-error" 
header "HX-Reswap" == "innerHTML show:#locale-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The locale is invalid." 

# Editing a locale without key shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
locale: en
key: 
value: hop
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The key is invalid." 

# Editing a locale without value shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
locale: en
key: hip
value: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#value-error" 
header "HX-Reswap" == "innerHTML show:#value-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The value is invalid." 

# Editing with correct data shows success
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
locale: en
key: hip
value: hop
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 
