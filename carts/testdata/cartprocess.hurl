# The cart list returns empty result
GET {{host}}/cart
HTTP 200
[Asserts]
body contains "The cart is empty"

# Add a product to the cart from the product detail page creates and cookie
POST {{host}}/cart/PDT1/add
HX-Current-Url: /product/t-shirt-tester-c-est-douter
HX-Request: true
[FormParams]
quantity: 1
HTTP 200
[Asserts]
cookie "cid" exists
cookie "cid[HttpOnly]" exists
cookie "cid[Max-Age]" exists
cookie "cid[Max-Age]" exists
cookie "cid[SameSite]" == "Strict"
cookie "cid[Path]" == "/"

# The cart list returns the product added
GET {{host}}/cart
HTTP 200
[Asserts]
body not contains "The cart is empty"

# Send delivery with correct redirect to payment page 
POST {{host}}/delivery
HX-Request: true
[FormParams]
delivery: colissimo
HTTP 200
[Asserts]
header "HX-Redirect" == "/cart/address"

# Saving the address works with data
POST {{host}}/cart/address
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
header "HX-Redirect" == "/payment"

# Send delivery with correct redirect to payment page 
POST {{host}}/payment
HX-Request: true
[FormParams]
payment: cash
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 
