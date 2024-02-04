# Delivery page redirects to cart when no product 
GET {{host}}/delivery.html
User-Agent: {{ua}}
HTTP 302
[Asserts]
header "Location" == "/cart.html"

# Add a product to the cart 
POST {{host}}/cart/PDT1/add.html
HX-Current-Url: /product/t-shirt-tester-c-est-douter.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
quantity: 1
HTTP 200

# Delivery page show delivery option 
GET {{host}}/delivery.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//input[@name=\"delivery\"]" exists

# Send delivery with empty value returns error
POST {{host}}/delivery.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
delivery: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#alert" 
header "HX-Reswap" == "innerHTML show:#alert:top" 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The delivery is invalid." 

# Send delivery with bad value returns error
POST {{host}}/delivery.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
delivery: abc
HTTP 200
[Asserts]
header "HX-Retarget" == "#alert" 
header "HX-Reswap" == "innerHTML show:#alert:top" 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The delivery is invalid." 

# Send delivery with correct redirect to payment page 
POST {{host}}/delivery.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
delivery: colissimo
HTTP 200
[Asserts]
header "HX-Redirect" == "/payment.html"