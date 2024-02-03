# Editing an order without status shows an error 
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status:
HTTP 200
[Asserts]
header "HX-Retarget" == "#status-error" 
header "HX-Reswap" == "innerHTML show:#status-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing an order with invalid status shows an error 
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status: hello
HTTP 200
[Asserts]
header "HX-Retarget" == "#status-error" 
header "HX-Reswap" == "innerHTML show:#status-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing an order with status shows success 
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status: delivering
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-success box']" exists 
xpath "normalize-space(//p/text())" == "The data has been saved successfully." 

# Accessing the order show the status selected
GET {{host}}/admin/orders/ORD1/edit.html
HTTP 200
[Asserts]
xpath "//option[@value='delivering'][@selected]" exists 

# Editing an order with status for the next iteration
POST {{host}}/admin/orders/ORD1/status.html
HX-Request: true
[FormParams]
status: processing
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-success box']" exists 