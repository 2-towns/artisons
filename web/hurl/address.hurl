# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: arnaud@yandex.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists

# Redirect to account when the otp is correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/otp.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email: arnaud@yandex.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/account/index.html"

# The address form is displayed
GET {{host}}/account/address.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//input[@name='firstname']" exists
xpath "//input[@name='lastname']" exists
xpath "//input[@name='street']" exists
xpath "//input[@name='complementary']" exists
xpath "//input[@name='city']" exists
xpath "//input[@name='zipcode']" exists
xpath "//input[@name='phone']" exists

# Saving the address without firstname shows an error 
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: 
lastname: Deville
street: 17 rue du moulin 
complementary: faubourg de la poste 
city: Lille 
zipcode: 59000
phone: 6012432122
HTTP 200
[Asserts]
header "HX-Retarget" == "#firstname-error" 
header "HX-Reswap" == "innerHTML show:#firstname-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Saving the address without lastname shows an error 
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: Arnaud
lastname: 
street: 17 rue du moulin 
complementary: faubourg de la poste 
city: Lille 
zipcode: 59000
phone: 6012432122
HTTP 200
[Asserts]
header "HX-Retarget" == "#lastname-error" 
header "HX-Reswap" == "innerHTML show:#lastname-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Saving the address without street shows an error 
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: Arnaud
lastname: Deville
street: 
complementary: faubourg de la poste 
city: Lille 
zipcode: 59000
phone: 6012432122
HTTP 200
[Asserts]
header "HX-Retarget" == "#street-error" 
header "HX-Reswap" == "innerHTML show:#street-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Saving the address without city shows an error 
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: Arnaud
lastname: Deville
street: 17 rue du moulin 
complementary: faubourg de la poste 
city:  
zipcode: 59000
phone: 6012432122
HTTP 200
[Asserts]
header "HX-Retarget" == "#city-error" 
header "HX-Reswap" == "innerHTML show:#city-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Saving the address without zipcode shows an error 
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: Arnaud
lastname: Deville
street: 17 rue du moulin 
complementary: faubourg de la poste 
city: Lille
zipcode: 
phone: 6012432122
HTTP 200
[Asserts]
header "HX-Retarget" == "#zipcode-error" 
header "HX-Reswap" == "innerHTML show:#zipcode-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Saving the address without phone shows an error 
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: Arnaud
lastname: Deville
street: 17 rue du moulin 
complementary: faubourg de la poste 
city: Lille
zipcode: 59000
phone: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#phone-error" 
header "HX-Reswap" == "innerHTML show:#phone-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Saving the address works with data
POST {{host}}/account/address.html
HX-Request: true
[FormParams]
firstname: Arnaud
lastname: Deville
street: 17 rue du moulin 
complementary: faubourg de la poste 
city: Lille
zipcode: 59000
phone: 6012432122
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 