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

# Adding the tags without key shows an error 
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key:  
label: Books
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the tags with bad key shows an error 
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: books!
label: Books
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the tags without Label key shows an error 
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: books
label: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#label-error" 
header "HX-Reswap" == "innerHTML show:#label-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the tags with bad score shows an error 
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: books
label: Books
score: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#score-error" 
header "HX-Reswap" == "innerHTML show:#score-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the tags with existing key shows an error 
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: mens
label: Books
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The tag exists already." 

# Adding the tags with root children key shows an error 
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: books3{{time}}
label: Books
children: womens
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The children cannot be root tags." 

# Adding the tags with correct data works
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: books{{time}}
label: Books
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/tags.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"

# Adding the tags with complete data works
POST {{host}}/admin/tags/add.html
HX-Request: true
[MultipartFormData]
key: books2{{time}}
label: Books
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
