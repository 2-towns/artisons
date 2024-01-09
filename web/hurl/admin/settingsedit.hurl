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

# Editing a product without name shows an error 
POST {{host}}/admin/settings.html
HX-Request: true
[MultipartFormData]
Name:
Email: contact@supershop.com
Address: 18 street mysupershop.com 
City: Oran
Zipcode: 31000
logo:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#name-error" 
header "HX-Reswap" == "innerHTML show:#name-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The name is invalid." 

# Editing a product without email shows an error 
POST {{host}}/admin/settings.html
HX-Request: true
[MultipartFormData]
Name: My shop
Email: 
Address: 18 street mysupershop.com 
City: Oran
Zipcode: 31000
logo:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#email-error" 
header "HX-Reswap" == "innerHTML show:#email-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The email is invalid." 