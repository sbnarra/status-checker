const servicesData = {};
let lastApiCallTime = new Date(new Date().getTime() - 3600000).toISOString().split('.')[0];

function formatDateTime(dateTimeStr) {
    const date = new Date(dateTimeStr);
    return date.toISOString().split('.')[0];
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function createTooltipContent(item) {
    let content = `Started: ${formatDateTime(item.started)}\nCompleted: ${formatDateTime(item.completed)}\n\n`;

    content += `Check Output:\n${item.check_output || 'None'}\n`;
    if (item.check_error) {
        content += `Check Error:\n${item.check_error}\n`;
    }

    if (item.recover_output) {
        content += `Recover Output:\n${item.recover_output}\n`;
    }
    if (item.recover_error) {
        content += `Recover Error:\n${item.recover_error}\n`;
    }

    if (item.recheck_output) {
        content += `Recheck Output:\n${item.recheck_output}\n`;
    }
    if (item.recheck_error) {
        content += `Recheck Error:\n${item.recheck_error}\n`;
    }

    return content;
}

function getStatusClass(status) {
    if (status === 'Success') {
        return 'status-success';
    } else if (status === 'Recovered') {
        return 'status-recovered';
    } else {
        return 'status-failed';
    }
}

function updateServices(data) {
    for (const [serviceName, historyItems] of Object.entries(data)) {
        if (!servicesData[serviceName]) {
            servicesData[serviceName] = [];
            const serviceDiv = document.createElement('div');
            serviceDiv.className = 'service';
            serviceDiv.id = `service-${serviceName}`;

            const serviceNameDiv = document.createElement('div');
            serviceNameDiv.className = 'service-name';
            serviceNameDiv.textContent = serviceName;
            serviceDiv.appendChild(serviceNameDiv);

            const statusBlocksDiv = document.createElement('div');
            statusBlocksDiv.className = 'status-blocks';
            serviceDiv.appendChild(statusBlocksDiv);

            document.getElementById('services-container').appendChild(serviceDiv);
        }

        const statusBlocksDiv = document.querySelector(`#service-${serviceName} .status-blocks`);

        historyItems.forEach(item => {
            servicesData[serviceName].push(item);

            const statusBlock = document.createElement('div');
            statusBlock.className = `status-block ${getStatusClass(item.status)}`;

            const tooltip = document.createElement('div');
            tooltip.className = 'tooltip';
            tooltip.textContent = createTooltipContent(item);

            statusBlock.appendChild(tooltip);
            statusBlocksDiv.appendChild(statusBlock);
        });
    }
}

async function fetchData() {
    try {
        const url = lastApiCallTime
            ? `/history?since=${encodeURIComponent(lastApiCallTime)}`
            : '/history';
        const response = await fetch(url);

        if (response.ok) {
            const data = await response.json();
            updateServices(data);
            lastApiCallTime = new Date().toISOString().split('.')[0];
        } else {
            console.error('Failed to fetch data:', response.statusText);
        }
    } catch (error) {
        console.error('Error fetching data:', error);
    }
}

// Initial fetch without 'since' to load all data
lastApiCallTime = '';
fetchData();

// Fetch new data every 30 seconds
setInterval(fetchData, 30000);