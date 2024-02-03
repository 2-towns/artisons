# Add a note without note shows an error 
POST {{host}}/admin/orders/ORD1/note.html
HX-Request: true
[FormParams]
note:
HTTP 200
[Asserts]
header "HX-Retarget" == "#note-error" 
header "HX-Reswap" == "innerHTML show:#note-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Add a note shows success 
POST {{host}}/admin/orders/ORD1/note.html
HX-Request: true
[FormParams]
note: Ta commande tu te la gardes !
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-success box']" exists 
xpath "normalize-space(//p/text())" == "The data has been saved successfully." 
xpath "normalize-space(//div[@class='box list-item'][last()]/b[text()='Ta commande tu te la gardes !'])" exists

# Accessing the order show the note 
GET {{host}}/admin/orders/ORD1/edit.html
HTTP 200
[Asserts]
xpath "normalize-space(//div[@class='box list-item'][last()]/b[text()='Ta commande tu te la gardes !'])" exists

