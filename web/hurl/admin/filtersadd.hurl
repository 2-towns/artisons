# Adding the filters without key shows an error 
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key:  
label: Materials
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the filters with bad key shows an error 
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key: materials!
label: Materials
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the filters without Label key shows an error 
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key: materials
label: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#label-error" 
header "HX-Reswap" == "innerHTML show:#label-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the filters with bad score shows an error 
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key: materials
label: Materials
score: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#score-error" 
header "HX-Reswap" == "innerHTML show:#score-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding the filters with existing key shows an error 
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key: colors
label: Materials
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The filter exists already." 

# Adding the filters with correct data works
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key: materials{{time}}
label: Materials
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/filters.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"

# Adding the filters with complete data works
POST {{host}}/admin/filters/add.html
HX-Request: true
[FormParams]
key: materials2{{time}}
label: Materials
children: materials;en
active: on
score: 10 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/filters.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
