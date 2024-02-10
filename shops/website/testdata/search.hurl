# The search returns products
GET {{host}}/search
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns products when min price is low
GET {{host}}/search?min=1
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product when min price is high
GET {{host}}/search?min=9999
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists

# The search returns products when max price is high
GET {{host}}/search?max=9999
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product when max price is low
GET {{host}}/search?max=1
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists

# The search returns products when keywords are found 
GET {{host}}/search?q=mug
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product whenkeywords are not found 
GET {{host}}/search?q=crazy
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists

# The search returns products when meta are found 
GET {{host}}/search?color=blue
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product when meta are not found 
GET {{host}}/search?color=crazy
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists





