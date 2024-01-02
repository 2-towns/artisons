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

# Add a note without note shows an error 
POST {{host}}/admin/orders/ORD1/note.html
HX-Request: true
[FormParams]
note:
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The note is required." 

# Add a note shows success 
POST {{host}}/admin/orders/ORD1/note.html
HX-Request: true
[FormParams]
note: Ta commande tu te la gardes !
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-success box']" exists 
xpath "normalize-space(//p/text())" == "The note has been added successfully." 
xpath "normalize-space(//div[@class='box list-item'][last()]/b[text()='Ta commande tu te la gardes !'])" exists

# Accessing the order show the note 
GET {{host}}/admin/orders/ORD1/edit.html
HTTP 200
[Asserts]
xpath "normalize-space(//div[@class='box list-item'][last()]/b[text()='Ta commande tu te la gardes !'])" exists

