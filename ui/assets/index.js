let sinceInput = document.getElementById('since-input')
let untilInput = document.getElementById('until-input')
let untilNowInput = document.getElementById('until-now-input')
let checksHistoryContainer = document.getElementById('checks-history-container')

let loadedChecksHistorys = {};
let historyLastPolled = apiDateTime(new Date(Date.now() - 3600 * 1000));

function apiDateTime(dateTime) {
  return dateTime.toISOString().split('.')[0];
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

function createHistoryTooltipContent(result) {
  let content = `<b>Status:</b> <i>${result.status}</i>
<b>Started:</b> <i>${new Date(result.started).toLocaleString()}</i>
<b>Completed:</b> <i>${new Date(result.completed).toLocaleString()}</i>
${createHistoryTooltipSection("Check", result.command, result.check_output, result.check_error)}`;
  let recover = createHistoryTooltipSection("Recover", result.recover, result.recover_output, result.recover_error);
  let recheck = createHistoryTooltipSection("Recheck", result.command, result.recheck_output, result.recheck_error);
  return content + (recover ? `${recover}` : "") + (recheck ? `${recheck}` : "");
}

function createHistoryTooltipSection(stage, command, output, error) {
  let content = ""
  if (output !== undefined) {
    content += `<b>${stage} Command:</b> <code>${command}</code>`;
  }
  if (error !== undefined) {
    content += `
<b>${stage} Error:</b> <code>${error}</code>`;
  }
  if (output !== undefined) {
    content += `
<b>${stage} Output:</b>
<code>${output || '""'}</code>`;
  }
  return content ? `<hr>${content}` : content
}

function getStatusClass(status) {
  return status === 'Success'   ? 'status-success'   :
         status === 'Recovered' ? 'status-recovered' :
                                  'status-failed'
}

function createHistoryDiv(name) {
  const historyDiv = document.createElement('div');
  historyDiv.className = 'check-history';
  historyDiv.id = `check-history-${name}`;

  const checkNameDiv = document.createElement('div');
  checkNameDiv.className = 'check-name';
  checkNameDiv.textContent = name;
  historyDiv.appendChild(checkNameDiv);

  const statusBlocksDiv = document.createElement('div');
  statusBlocksDiv.className = 'status-blocks';
  historyDiv.appendChild(statusBlocksDiv);

  loadedChecksHistorys[name] = [];
  checksHistoryContainer.appendChild(historyDiv);
}

function updateHistoryDiv(name, history) {
  const statusBlocksDiv = document.querySelector(`#check-history-${name} .status-blocks`);
  history.forEach(result => {
    loadedChecksHistorys[name].push(result);

    const statusBlock = document.createElement('div');
    statusBlock.className = `status-block ${getStatusClass(result.status)}`;

    const tooltip = document.createElement('div');
    tooltip.className = 'history-tooltip';
    tooltip.innerHTML = createHistoryTooltipContent(result);
    statusBlock.appendChild(tooltip);

    statusBlock.addEventListener('mousemove', function(e) {
      if (e.clientX < (window.innerWidth / 2)) {
        tooltip.style.left = '100%';
        tooltip.style.right = 'auto';
      } else {
        tooltip.style.left = 'auto';
        tooltip.style.right = '100%';
      }
    });

    statusBlocksDiv.appendChild(statusBlock);
  });
}

function updateHistorys(historys) {
  for (const [name, history] of Object.entries(historys)) {
    if (!loadedChecksHistorys[name]) {
      createHistoryDiv(name)
    }
    updateHistoryDiv(name, history)
  }
}

async function fetchChecksHistorys() {
  try {
    const url = untilNowInput.checked
      ? `${window.location.origin}/history?since=${historyLastPolled}`
      : `${window.location.origin}/history?since=${sinceInput.value}&until=${untilInput.value}`;
    let newHistoryLastPolled = apiDateTime(new Date());
    const response = await fetch(url);
    if (response.ok) {
      const historys = await response.json();
      updateHistorys(historys);
      historyLastPolled = newHistoryLastPolled
    } else {
      console.error('Failed to fetch data:', response.statusText);
    }
  } catch (error) {
    console.error('Error fetching data:', error);
  }
}

function resetPage() {
  loadedChecksHistorys = {}
  checksHistoryContainer.innerHTML = ""
  historyLastPolled = sinceInput.value;
  fetchChecksHistorys();
}

sinceInput.value = historyLastPolled
untilInput.value = apiDateTime(new Date())
untilInput.disabled = true
untilNowInput.checked = true

fetchChecksHistorys();
let fetchChecksHistorysInterval = 10_000
let fetchChecksHistorysId = setInterval(fetchChecksHistorys, fetchChecksHistorysInterval);

untilNowInput.addEventListener('click', () => {
  untilInput.disabled = untilNowInput.checked
  untilInput.value = apiDateTime(new Date())
  if (untilNowInput.checked) {
    resetPage()
    fetchChecksHistorysId = setInterval(fetchChecksHistorys, fetchChecksHistorysInterval);
  } else {
    clearInterval(fetchChecksHistorysId);
  }
});
document.getElementById('filter-button').addEventListener('click', resetPage);
