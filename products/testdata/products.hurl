# Display the products
GET {{host}}/admin/products
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//input[@name='q'][@hx-get='/admin/products?page=1']" exists
xpath "//table/tbody/tr" count > 0

# Access the products with huge page shoes empty message
GET {{host}}/admin/products?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

# Search products by id 
GET {{host}}/admin/products?q=PDT1
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//td[text()='PDT1']" exists
xpath "//a[@href='/admin/products/PDT1/edit']" exists
xpath "//a[@hx-post='/admin/products/PDT1/delete']" exists

# Search products by sku 
GET {{host}}/admin/products?q=SKU1
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//td[text()='PDT1']" exists
xpath "//a[@href='/admin/products/PDT1/edit']" exists
xpath "//a[@hx-post='/admin/products/PDT1/delete']" exists

# Search products by title 
GET {{host}}/admin/products?q=tester
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//td[text()='PDT1']" exists
xpath "//a[@href='/admin/products/PDT1/edit']" exists
xpath "//a[@hx-post='/admin/products/PDT1/delete']" exists

# Search products by non existing keyword shows empty message 
GET {{host}}/admin/products?q=totoisverydangerous
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

