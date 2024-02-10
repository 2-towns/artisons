# Editing the filters without label key shows an error 
POST {{host}}/admin/filters/sizes/edit
HX-Request: true
[FormParams]
key: sizes
label: 
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
header "HX-Retarget" == "#label-error" 
header "HX-Reswap" == "innerHTML show:#label-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the filters with bad score shows an error 
POST {{host}}/admin/filters/sizes/edit
HX-Request: true
[FormParams]
key: sizes
label: sizes
score: hello
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
header "HX-Retarget" == "#score-error" 
header "HX-Reswap" == "innerHTML show:#score-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the filters with correct data works
POST {{host}}/admin/filters/sizes/edit
HX-Request: true
[FormParams]
key: sizes
label: sizes
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
header "HX-Redirect" == "/admin/filters"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"

# Editing the filters with complete data works
POST {{host}}/admin/filters/sizes/edit
HX-Request: true
[FormParams]
key: sizes
label: sizes
children: books;en
root: on
score: 10 
image:file,image.jpg;
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
header "HX-Redirect" == "/admin/filters"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
