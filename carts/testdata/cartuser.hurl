# The cart list returns empty result
GET {{host}}/cart
[Cookies]
wsid: 333333 
HTTP 200
[Asserts]
body contains "The cart is empty"

# Add a product to the cart from the product detail page
POST {{host}}/cart/PDT1/add
HX-Current-Url: /product/t-shirt-tester-c-est-douter
HX-Request: true
[FormParams]
quantity: 3
HTTP 200

# The cart list returns the product added
GET {{host}}/cart
HTTP 200
[Asserts]
body not contains "The cart is empty."
cookie "cid" exists

# Delete few quantity keeps the product into the cart
POST {{host}}/cart/PDT1/delete
HX-Current-Url: /product/t-shirt-tester-c-est-douter
HX-Request: true
[FormParams]
quantity: 1
HTTP 200

# The cart list returns the product added
GET {{host}}/cart
HTTP 200
[Asserts]
body not contains "The cart is empty"
cookie "cid" exists

# Delete a product to the cart from the empty page
POST {{host}}/cart/PDT1/delete
HX-Current-Url: /product/t-shirt-tester-c-est-douter
HX-Request: true
[FormParams]
quantity: 99
HTTP 200

# The cart list returns empty result
GET {{host}}/cart
HTTP 200
[Asserts]
body contains "The cart is empty"
