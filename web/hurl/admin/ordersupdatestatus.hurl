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

# Editing an order without status shows an error 
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status:
HTTP 200
[Asserts]
header "HX-Retarget" == "#status-error" 
header "HX-Reswap" == "innerHTML show:#status-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The status is invalid." 

# Editing an order with invalid status shows an error 
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#status-error" 
header "HX-Reswap" == "innerHTML show:#status-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The status is invalid." 

# Editing an order with status shows success 
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status: delivering
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-success box']" exists 
xpath "normalize-space(//p/text())" == "The status has been updated successfully." 

# Accessing the order show the status selected
GET {{host}}/admin/orders/ORD1/edit.html
HTTP 200
[Asserts]
xpath "//option[@value='delivering'][@selected]" exists 

# Editing an order with status for the next iteration
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status: processing
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-success box']" exists 