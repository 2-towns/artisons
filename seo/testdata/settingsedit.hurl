# Editing the contact settings without name shows an error 
POST {{host}}/admin/contact-settings
HX-Request: true
[Cookies]
wsid: 444444 
[MultipartFormData]
name:
email: contact@supershop.com
address: 18 street mysupershop.com 
city: Oran
zipcode: 31000
logo:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#name-error" 
header "HX-Reswap" == "innerHTML show:#name-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the contact settings with invalid email shows an error 
POST {{host}}/admin/contact-settings
HX-Request: true
[Cookies]
wsid: 444444 
[MultipartFormData]
name: My shop
email: 
address: 18 street mysupershop.com 
city: Oran
zipcode: 31000
phone: 12345
logo:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#email-error" 
header "HX-Reswap" == "innerHTML show:#email-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the contact settings without phone shows an error 
POST {{host}}/admin/contact-settings
HX-Request: true
[Cookies]
wsid: 444444 
[MultipartFormData]
name: My shop
email: contact@shop.com
address: 18 street mysupershop.com 
city: Oran
phone: 
logo:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#phone-error" 
header "HX-Reswap" == "innerHTML show:#phone-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the contact settings with correct data shows success
POST {{host}}/admin/contact-settings
HX-Request: true
[Cookies]
wsid: 444444 
[MultipartFormData]
name: My shop
email: contact@shop.com
address: 18 street mysupershop.com 
city: Oran
zipcode: 31000
phone: 123456
logo:file,image.jpg;
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 

# Editing the shop settings with invalid items shows an error 
POST {{host}}/admin/shop-settings
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
items: hello
min: 10 
last: 10 
active: on
HTTP 200
[Asserts]
header "HX-Retarget" == "#items-error" 
header "HX-Reswap" == "innerHTML show:#items-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the shop settings with invalid min shows an error 
POST {{host}}/admin/shop-settings
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
items: 10
min: hello
last: 10 
active: on
HTTP 200
[Asserts]
header "HX-Retarget" == "#min-error" 
header "HX-Reswap" == "innerHTML show:#min-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the shop settings with invalid image width shows an error 
POST {{host}}/admin/shop-settings
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
items: 10
min: 10
last: 10 
active: on
image_width: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#image_width-error" 
header "HX-Reswap" == "innerHTML show:#image_width-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the shop settings with invalid image height shows an error 
POST {{host}}/admin/shop-settings
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
items: 10
min: 10
last: 10 
active: on
image_height: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#image_height-error" 
header "HX-Reswap" == "innerHTML show:#image_height-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the shop settings shows success 
POST {{host}}/admin/shop-settings
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
items: 10
min: 10
last: 10 
active: on
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 