$(document).ready(function() {
    // Init Variables
    let colors = {
        pattern: ['#60B044', '#F6C600', '#F97600', '#FF0000'],
        threshold: {
            values: [50, 70, 90, 100]
        }
    };
    let gaugeConfig = {
        columns: [
            ['usage', 0.0]
        ],
        type: 'gauge',
    };

    // CPU Chart
    let cpuChart = c3.generate({
        bindto: '#cpu',
        data: gaugeConfig,
        color: colors,
        size: {
            height: 300
        }
    });

    let memChart = c3.generate({
        bindto: '#memory',
        data: gaugeConfig,
        color: colors,
        size: {
            height: 300
        }
    });

    let loadChart = c3.generate({
        bindto: '#load',
        data: {
            columns: [
                ['load1', 0.0],
                ['load5', 0.0],
                ['load15', 0.0]
            ],
            type: 'gauge',
        },
        color: colors,
        size: {
            height: 300
        }
    });

    // Websocket config
    url = 'ws://' + window.location.host + '/ws';
    c = new WebSocket(url);


    // Handle WS message
    c.onmessage = function (msg) {
        let data = JSON.parse(msg.data);

        cpuChart.load({columns: [["usage", data.cpu.perc]]});
        memChart.load({columns: [["usage", data.memory.perc]]});
        loadChart.load({columns: [
            ["load1", data.load.load_1 * 100],
            ["load5", data.load.load_5 * 100],
            ["load15", data.load.load_15 * 100]
        ]});

        // Sort disks
        let diskData = data.disks
        diskData.sort(function(a, b) {
            if (a.mount == b.mount) return 0;
            if (a.mount > b.mount) return 1;
            if (a.mount < b.mount) return -1;
        });

        // Render table
        let table = $("#table-content");
        table.html("");
        data.disks.forEach(function(disk) {
            let tr = $("<tr>");
            tr.append("<td>" + disk.mount + "</td>");
            tr.append("<td>" + disk.used + "</td>");
            tr.append("<td>" + disk.free + "</td>");
            tr.append("<td>" + disk.perc + "%</td>");

            if (disk.perc <= 50) {
                tr.append("<td class='text-center'><span class='badge' style='background-color: #60B044'>&nbsp;&nbsp;&nbsp;</td>");
            } else if (disk.perc <= 70) {
                tr.append("<td class='text-center'><span class='badge' style='background-color: #F6C600'>&nbsp;&nbsp;&nbsp;</td>");
            } else if (disk.perc <= 70) {
                tr.append("<td class='text-center'><span class='badge' style='background-color: #F97600'>&nbsp;&nbsp;&nbsp;</td>");
            } else {
                tr.append("<td class='text-center'><span class='badge' style='background-color: #FF0000'>&nbsp;&nbsp;&nbsp;</td>");
            }

            table.append(tr);
        })
    };
});
