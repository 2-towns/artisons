# Delivery page redirects to cart when no product 
GET {{host}}/delivery
HTTP 401

# Add a product to the cart 
POST {{host}}/cart/PDT1/add
HX-Current-Url: /product/t-shirt-tester-c-est-douter
HX-Request: true
[FormParams]
quantity: 1
HTTP 200

# Delivery page show delivery option 
GET {{host}}/delivery
HTTP 200
[Asserts]
xpath "//input[@name=\"delivery\"]" exists

# Send delivery with empty value returns error
POST {{host}}/delivery
HX-Request: true
[FormParams]
delivery: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#alert" 
header "HX-Reswap" == "innerHTML show:#alert:top" 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The delivery is invalid." 

# Send delivery with bad value returns error
POST {{host}}/delivery
HX-Request: true
[FormParams]
delivery: abc
HTTP 200
[Asserts]
header "HX-Retarget" == "#alert" 
header "HX-Reswap" == "innerHTML show:#alert:top" 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The delivery is invalid." 

# Send delivery with correct redirect to payment page 
POST {{host}}/delivery
HX-Request: true
[FormParams]
delivery: colissimo
HTTP 200
[Asserts]
header "HX-Redirect" == "/cart/address"