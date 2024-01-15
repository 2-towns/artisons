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

# Display the orders
GET {{host}}/admin/orders.html
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//input[@name='q'][@hx-get='/admin/orders.html?page=1']" exists
xpath "//table/tbody/tr" count > 0

# Access the orders with huge page shoes empty message
GET {{host}}/admin/orders.html?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

# Search orders by id 
GET {{host}}/admin/orders.html?q=ORD1
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "string((//td)[1])" == "ORD1"
xpath "//a[@href='/admin/orders/ORD1/edit.html']" exists

# Search orders by status 
GET {{host}}/admin/orders.html?q=created
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found."
xpath "normalize-space(//td[text()='ORD2'])" exists
xpath "//a[@href='/admin/orders/ORD2/edit.html']" exists

# Search orders by delivery 
GET {{host}}/admin/orders.html?q=Home
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "normalize-space(//td[text()='ORD1'])" exists
xpath "//a[@href='/admin/orders/ORD1/edit.html']" exists

# Search orders by payment 
GET {{host}}/admin/orders.html?q=Card
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "normalize-space(//td[text()='ORD1'])" exists
xpath "//a[@href='/admin/orders/ORD1/edit.html']" exists

# Search orders by non existing keyword shows empty message 
GET {{host}}/admin/orders.html?q=totoisverydangerous
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

