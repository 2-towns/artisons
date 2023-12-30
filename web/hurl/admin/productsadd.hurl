# Display otp form when the correct email is used
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/auth/login.html']" exists
[Captures]
glue: xpath "normalize-space(//input[@id='glue']/@value)"

# Redirect to dashboard when the otp is correct
POST {{host}}/auth/login.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
glue: {{glue}}
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Adding a product without price shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price:
quantity: 1
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The price is invalid." 

# Adding a product with an invalid price shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: abc
quantity: 1
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The price is invalid." 

# Adding a product without title shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: 
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The title is invalid." 

# Adding a product without description shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 
sku: 123
price: 123.5
quantity: 1
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The description is invalid." 

# Adding a product with bad sku shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123!
price: 123.5
quantity: 1
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The sku is invalid." 

# Adding a product without quantity shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The quantity is invalid." 

# Adding a product with bad quantity shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: abc
status: online
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The quantity is invalid." 

# Adding a product without status shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: 
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The status is invalid." 

# Adding a product with bad status shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: abc
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The status is invalid." 

# Adding a product with bad weight shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
weight: abc
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The weight is invalid." 

# Adding a product with bad discount shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
discount: abc
tags: winter cold
image_1:file,image.jpg;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The discount is invalid." 

# Adding a product without file shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
tags: winter cold
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The image is required." 

# Adding a product with a bad picture shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
tags: winter cold
image_1:file,image.txt;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The image is invalid." 

# Adding a product with a bad image shows an error 
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
tags: winter cold
image_1:file,image.jpg;
image_2:file,image.txt;
HTTP 200
[Asserts]
xpath "//div[@class='form-error']" exists 
xpath "normalize-space(//div[@class='form-error']/text())" == "The image is invalid." 

# Adding a product with correct data works
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
sku: 123
price: 123.5
quantity: 1
status: online
weight: 12.4
discount: 12.5
tags: winter cold
image_1:file,image.jpg;
image_2:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/products.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"

# Adding a product with only required fields works
POST {{host}}/admin/products/add.html
HX-Request: true
[MultipartFormData]
title: T-shirt développeur unisexe JavaScript Park
description: 100 % coton pour les couleurs unies
price: 123.5
quantity: 1
status: online
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/products.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"