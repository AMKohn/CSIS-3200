(function() {
    const loadTemplates = function () {
        let templates = {
            last60: "last-60-mins",
            slowestWeb: "slowest-web",
            stats: "stats",
            hosts: "hosts"
        };

        for (let tName in templates) {
            templates[tName] = {
                template: Hogan.compile(document.getElementById(templates[tName]).innerHTML),
                container: document.querySelector("." + templates[tName])
            };
        }

        return templates;
    };

    const parseData = function(d) {
        d.stats.webRequests = numeral(d.stats.webRequests).format("0.0a");
        d.stats.databaseQueries = numeral(d.stats.databaseQueries).format("0.0a");
        d.stats.searchQueries = numeral(d.stats.searchQueries).format("0.0a");
        d.stats.messageRate = numeral(d.stats.messageRate).format("0.0a");

        d.stats.cacheHitRate += "%";
        d.stats.cpuUsage += "%";
        d.stats.overallErrorRate += "%";
        d.stats.webResponseTime += "ms";

        d.hosts.forEach(e => {
            e.errorRate += "%";
            e.cpuUsage += "%";
            e.responseTime += " ms";

            e.memory = numeral(e.memoryUsage * 1000).format("0.0 ib") + " / " + numeral(e.memoryCapacity * 1000).format("0.0 ib");
            e.throughput = numeral(e.throughput).format("0.0a") + " RPM";
        });

        return d;
    };

    let chart = null;

    const renderChart = function(data) {
        if (chart === null) {
            chart = new Chart("resp-time", {
                type: "line",
                data: {
                    labels: ["-30 mins", "-25 mins", "-20 mins", "-15 mins", "-10 mins", "-5 mins", "0 mins"],
                    datasets: [{
                        data: data.responseTimes,
                        borderColor: "#ff3b30",
                        backgroundColor: "rgba(255, 59, 48, 0.2)"
                    }]
                },
                options: {
                    legend: false,
                    tooltips: true,
                    maintainAspectRatio: false
                }
            });
        }
        else {
            chart.data.datasets.forEach((dataset) => {
                dataset.data = data.responseTimes;
            });

            chart.update();
        }
    };

    const renderTemplates = function(data) {
        for (let t in templates) {
            templates[t].container.innerHTML = templates[t].template.render(data);
        }
    };

    window.templates = loadTemplates();

    if (typeof window.initData !== "undefined") {
        window.initData = parseData(window.initData);
    }

    renderChart(window.initData);

    renderTemplates(window.initData);

    const updateData = () => {
        fetch("/api/getData")
            .then(response => response.json())
            .then(d => {
                d = parseData(d);

                renderChart(d);

                renderTemplates(d);
            });
    };

    // Update every 30 seconds
    setInterval(updateData, 30 * 1000);

    updateData();
})();