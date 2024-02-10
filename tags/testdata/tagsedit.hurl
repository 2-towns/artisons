# Editing the tags without label key shows an error 
POST {{host}}/admin/tags/mens/edit
HX-Request: true
[MultipartFormData]
key: mens
label: 
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
header "HX-Retarget" == "#label-error" 
header "HX-Reswap" == "innerHTML show:#label-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the tags with bad score shows an error 
POST {{host}}/admin/tags/mens/edit
HX-Request: true
[MultipartFormData]
key: mens
label: Mens
score: hello
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
header "HX-Retarget" == "#score-error" 
header "HX-Reswap" == "innerHTML show:#score-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the tags with root children key shows an error 
POST {{host}}/admin/tags/mens/edit
HX-Request: true
[MultipartFormData]
key: mens
label: Books
children: womens
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The children cannot be root tags." 

# Editing the tags with correct data works
POST {{host}}/admin/tags/mens/edit
HX-Request: true
[MultipartFormData]
key: mens
label: Mens
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/tags"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"

# Editing the tags with complete data works
POST {{host}}/admin/tags/mens/edit
HX-Request: true
[MultipartFormData]
key: mens
label: Mens
children: books;en
root: on
score: 10 
image:file,image.jpg;
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/tags"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
