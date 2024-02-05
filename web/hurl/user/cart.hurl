# Delete a product to the cart from the empty page
POST {{host}}/cart/PDT1/delete.html
HX-Current-Url: /product/t-shirt-tester-c-est-douter.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
quantity: 99
HTTP 200

# Delete a product to the cart from the empty page
POST {{host}}/cart/PDT2/delete.html
HX-Current-Url: /product/t-shirt-tester-c-est-douter.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
quantity: 99
HTTP 200

# The cart list returns empty result
GET {{host}}/cart.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "The cart is empty"

# Add a product to the cart from the product detail page
POST {{host}}/cart/PDT1/add.html
HX-Current-Url: /product/t-shirt-tester-c-est-douter.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
quantity: 1
HTTP 200

# The cart list returns the product added
GET {{host}}/cart.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "The cart is empty"
cookie "cid" not exists

# Add a product to the cart from cart page returns the cart updated
POST {{host}}/cart/PDT1/add.html
HX-Current-Url: /cart.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
quantity: 1
HTTP 200
[Asserts]
xpath "//p[@class='product']" exists

# Delete a product to the cart from the empty page
POST {{host}}/cart/PDT1/delete.html
HX-Current-Url: /product/t-shirt-tester-c-est-douter.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
quantity: 99
HTTP 200

# The cart list returns empty result
GET {{host}}/cart.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "The cart is empty"
