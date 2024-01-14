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

# Display the products
GET {{host}}/admin/products.html
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//input[@name='q'][@hx-get='/admin/products.html?page=1']" exists
xpath "//table/tbody/tr" count > 0

# Access the products with huge page shoes empty message
GET {{host}}/admin/products.html?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

# Search products by id 
GET {{host}}/admin/products.html?q=PDT2
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "string((//td)[1])" == "PDT2"
xpath "//a[@href='/admin/products/PDT2/edit.html']" exists
xpath "//a[@hx-post='/admin/products/PDT2/delete.html']" exists

# Search products by sku 
GET {{host}}/admin/products.html?q=SKU3
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "string((//td)[1])" == "PDT3"
xpath "//a[@href='/admin/products/PDT3/edit.html']" exists
xpath "//a[@hx-post='/admin/products/PDT3/delete.html']" exists

# Search products by title 
GET {{host}}/admin/products.html?q=T-shirt%20d%C3%A9veloppeur%20unisexe%20JavaScript
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//td[text()='PDT2']" exists
xpath "//a[@href='/admin/products/PDT2/edit.html']" exists
xpath "//a[@hx-post='/admin/products/PDT2/delete.html']" exists

# Search products by non existing keyword shows empty message 
GET {{host}}/admin/products.html?q=totoisverydangerous
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

