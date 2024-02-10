# Display the orders
GET {{host}}/admin/orders
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//input[@name='q'][@hx-get='/admin/orders?page=1']" exists
xpath "//table/tbody/tr" count > 0

# Access the orders with huge page shoes empty message
GET {{host}}/admin/orders?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

# Search orders by id 
GET {{host}}/admin/orders?q=ORD1
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "string((//td)[1])" == "ORD1"
xpath "//a[@href='/admin/orders/ORD1/edit']" exists

# Search orders by status 
GET {{host}}/admin/orders?q=created
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found."
xpath "normalize-space(//td[text()='ORD2'])" exists
xpath "//a[@href='/admin/orders/ORD2/edit']" exists

# Search orders by delivery 
GET {{host}}/admin/orders?q=Home
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "normalize-space(//td[text()='ORD1'])" exists
xpath "//a[@href='/admin/orders/ORD1/edit']" exists

# Search orders by payment 
GET {{host}}/admin/orders?q=Card
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "normalize-space(//td[text()='ORD1'])" exists
xpath "//a[@href='/admin/orders/ORD1/edit']" exists

# Search orders by non existing keyword shows empty message 
GET {{host}}/admin/orders?q=totoisverydangerous
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

