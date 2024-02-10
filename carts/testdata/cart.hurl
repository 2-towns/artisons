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

# Add a product to the cart from the product detail page creates and cookie
POST {{host}}/cart/PDT1/delete
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

# The cart list returns empty result
GET {{host}}/cart
HTTP 200
[Asserts]
body contains "The cart is empty"
