const findById = document.getElementById.bind(document);
const findByClassName = document.getElementsByClassName.bind(document);

let now = new Date();
let days = 7;
let labels = [];
let visitsData = [];
let uniqueVisitsData = [];
let pageviewsData = [];
let bounceRatesData = [];
let ordersData = [];
let ordersCountsData = [];
let schart;
let ochart;
let visitsChartIndex = 1;
let ordersChartIndex = 1;

function createVisitsChart(data, labels) {
	schart = new Chartist.Line(
		"#stats-chart",
		{
			labels: labels,
			series: [data],
		},
		{
			low: 0,
			showArea: false,
			plugins: [Chartist.plugins.tooltip()],
		},
	);
}

function createOrdersChart(data, labels) {
	ochart = new Chartist.Line(
		"#orders-chart",
		{
			labels: labels,
			series: [data],
		},
		{
			low: 0,
			showArea: true,
			plugins: [Chartist.plugins.tooltip()],
		},
	);
}

function reload() {
	visitsData = JSON.parse(findById("visits").innerText);
	uniqueVisitsData = JSON.parse(findById("unique-visits").innerText);
	pageviewsData = JSON.parse(findById("pageviews").innerText);
	bounceRatesData = JSON.parse(findById("bounce-rates").innerText);
	ordersData = JSON.parse(findById("orders").innerText);
	ordersCountsData = JSON.parse(findById("orders-counts").innerText);
	days = parseInt(document.getElementById("days").value, 10);
	labels = [];
	now = new Date();

	for (let i = 0; i < days; i++) {
		const month = now.toLocaleString("default", { month: "short" });
		labels.unshift(month + " " + now.getDate());
		now.setDate(now.getDate() - 1);
	}

	switch (visitsChartIndex) {
		case 1: {
			toggleSelectedClass("visits-anchor");
			createVisitsChart(visitsData, labels);
			break;
		}

		case 2: {
			toggleSelectedClass("unique-visits-anchor");
			createVisitsChart(uniqueVisitsData, labels);
			break;
		}

		case 3: {
			toggleSelectedClass("pageviews-anchor");
			createVisitsChart(pageviewsData, labels);
			break;
		}

		case 4: {
			toggleSelectedClass("bounce-rates-anchor");
			createVisitsChart(bounceRatesData, labels);
			break;
		}
	}

	switch (ordersChartIndex) {
		case 1: {
			toggleSelectedClass("orders-anchor", "stats-orders");
			createOrdersChart(ordersData, labels);
			break;
		}
		case 2: {
			toggleSelectedClass("orders-counts-anchor", "stats-orders");
			createOrdersChart(ordersCountsData, labels);
			break;
		}
	}

	registerEvents();
}

reload();

document.body.addEventListener("ecm-dashboard-reload", function (evt) {
	reload();
});

function toggleSelectedClass(id, prefix = "stats") {
	const elts = document.getElementsByClassName(prefix + "-selected");

	for (const elt of elts) {
		elt.classList.remove(prefix + "-selected");
	}

	findById(id).classList.add(prefix + "-selected");
}

function registerEvents() {
	const visits = document.getElementById("visits-anchor");

	visits.onclick = function (e) {
		e.preventDefault();

		toggleSelectedClass("visits-anchor");

		schart.update({ series: [visitsData], labels });

		visitsChartIndex = 1;
	};

	const uniques = document.getElementById("unique-visits-anchor");

	uniques.onclick = function (e) {
		e.preventDefault();

		toggleSelectedClass("unique-visits-anchor");

		schart.update({ series: [uniqueVisitsData], labels });

		visitsChartIndex = 2;
	};

	const pageviews = document.getElementById("pageviews-anchor");

	pageviews.onclick = function (e) {
		e.preventDefault();

		toggleSelectedClass("pageviews-anchor");

		schart.update({ series: [pageviewsData], labels });

		visitsChartIndex = 3;
	};

	const bouceRates = document.getElementById("bounce-rates-anchor");

	bouceRates.onclick = function (e) {
		e.preventDefault();

		toggleSelectedClass("bounce-rates-anchor");

		schart.update({ series: [bounceRatesData], labels });

		visitsChartIndex = 4;
	};

	const orders = document.getElementById("orders-anchor");

	orders.onclick = function (e) {
		e.preventDefault();

		toggleSelectedClass("orders-anchor", "stats-orders");

		ochart.update({ series: [ordersData], labels });

		ordersChartIndex = 1;
	};

	const ordersCounts = document.getElementById("orders-counts-anchor");

	ordersCounts.onclick = function (e) {
		e.preventDefault();

		toggleSelectedClass("orders-counts-anchor", "stats-orders");

		ochart.update({ series: [ordersCountsData], labels });

		ordersChartIndex = 2;
	};
}
