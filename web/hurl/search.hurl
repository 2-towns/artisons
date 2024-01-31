# The search returns products
GET {{host}}/search.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns products when min price is low
GET {{host}}/search.html?min=1
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product when min price is high
GET {{host}}/search.html?min=9999
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists

# The search returns products when max price is high
GET {{host}}/search.html?max=9999
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product when max price is low
GET {{host}}/search.html?max=1
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists

# The search returns products when keywords are found 
GET {{host}}/search.html?q=mug
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product whenkeywords are not found 
GET {{host}}/search.html?q=crazy
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists

# The search returns products when meta are found 
GET {{host}}/search.html?color=blue
User-Agent: {{ua}}
HTTP 200
[Asserts]
body not contains "No results found."
xpath "//div[@class='product']" exists

# The search returns zero product when meta are not found 
GET {{host}}/search.html?color=crazy
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."
xpath "//div[@class='product']" not exists





