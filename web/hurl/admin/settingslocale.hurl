# Editing a locale without key shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
key: 
value: hop
HTTP 200
[Asserts]
header "HX-Retarget" == "#key-error" 
header "HX-Reswap" == "innerHTML show:#key-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing a locale without value shows an error 
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
key: hip
value: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#value-error" 
header "HX-Reswap" == "innerHTML show:#value-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing with correct data shows success
POST {{host}}/admin/locale.html
HX-Request: true
[FormParams]
key: hip
value: hop
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 
