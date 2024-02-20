# Editing the shop settings to redirect after adding to cart shows success 
POST {{host}}/admin/shop-settings
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
items: 10
min: 10
last: 10 
active: on
redirect: on
HTTP 200
[Asserts]
xpath "normalize-space(//b[@class='alert-success-title alert-title']/text())" == "Success" 

# Add a product to the cart from the product detail page creates and cookie
POST {{host}}/cart/PDT1/add
HX-Current-Url: /product/t-shirt-tester-c-est-douter
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
quantity: 1
HTTP 200
[Asserts]
header "HX-Redirect" == "/cart"

# Editing the shop settings to not redirect after adding to cart shows success 
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