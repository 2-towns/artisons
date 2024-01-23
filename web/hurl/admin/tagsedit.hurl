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

# Editing the tags without label key shows an error 
POST {{host}}/admin/tags/mens/edit.html
HX-Request: true
[MultipartFormData]
key: mens
label: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#label-error" 
header "HX-Reswap" == "innerHTML show:#label-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the tags with bad score shows an error 
POST {{host}}/admin/tags/mens/edit.html
HX-Request: true
[MultipartFormData]
key: mens
label: Mens
score: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#score-error" 
header "HX-Reswap" == "innerHTML show:#score-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the tags with root children key shows an error 
POST {{host}}/admin/tags/mens/edit.html
HX-Request: true
[MultipartFormData]
key: mens
label: Books
children: womens
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The children cannot be root tags." 

# Editing the tags with correct data works
POST {{host}}/admin/tags/mens/edit.html
HX-Request: true
[MultipartFormData]
key: mens
label: Mens
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/tags.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"

# Editing the tags with complete data works
POST {{host}}/admin/tags/mens/edit.html
HX-Request: true
[MultipartFormData]
key: mens
label: Mens
children: books;en
root: on
score: 10 
image:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/tags.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
