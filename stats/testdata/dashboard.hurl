# Display the dashboard
GET {{host}}/admin/index
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
xpath "//select[@name='days'][@hx-get='/admin/index']" exists
xpath "//script[@id='visits']" exists
xpath "//script[@id='unique-visits']" exists
xpath "//script[@id='pageviews']" exists
xpath "//script[@id='bounce-rates']" exists
xpath "//script[@id='orders']" exists
xpath "//script[@id='orders-counts']" exists
xpath "//script[@src='/js/admin/chartist.js']" exists
xpath "//script[@src='/js/admin/chartist-plugin-tooltip.js']" exists
xpath "//input[@id='days'][@value='7']" exists

# Change the pagination
GET {{host}}/admin/index?days=14
HX-Request: true
HTTP 200
[Asserts]
xpath "//select[@name='days'][@hx-get='/admin']" not exists
xpath "//script[@src='/js/admin/chartist.js']" not exists
xpath "//script[@src='/js/admin/chartist-plugin-tooltip.js']" not exists
xpath "//script[@id='visits']" exists
xpath "//script[@id='unique-visits']" exists
xpath "//script[@id='pageviews']" exists
xpath "//script[@id='bounce-rates']" exists
xpath "//script[@id='orders']" exists
xpath "//script[@id='orders-counts']" exists
xpath "//input[@id='days'][@value='14']" exists