const routes = {
    "^#home$": displayHome,
    "^#position$": displayPositionList,
    "^#position/new$": displaySaveNew,
    "^#position/level$": displayLevel,
    "^#position/([0-9]+)$": displayEditExiting,
    "^#position/([0-9]+)/recall$": displayPositionRecall,
    "^#setting/calibration": displayCalibration,
    "^#setting/preference$": displayPreferences,
};

function loadContent() {

    const div = getAppDiv();

    const url = location.hash;

    for (const key in routes) {
        const matches = url.match(key);
        if (matches) {
            console.log(url + " matches " + key);
            routes[key](matches, div);
            break;
        }
    }
}

if (!location.hash) {
    location.hash = "#home";
}

loadContent();

window.addEventListener("hashchange", loadContent)

function getAppDiv() {
    return document.getElementById("app")
}

async function getUrls(urls) {
    return await Promise.all(urls.map(url => fetch(url).then(r => r.json()).catch(err => console.error(err))))
}

async function getUrl(url) {
    const results = await getUrls([url])
    return results[0]
}

function putUrl(url, body) {
    console.log(`Updating ${url} with ${JSON.stringify(body)}`)
    return fetch(url, {
        method: 'PUT',
        body: JSON.stringify(body),
    })
    .then(r => r.json())
    .catch(err => {
        console.error(err)
        displayAlert('alert-danger', 'Update failed', false)
    });
}

function postUrl(url, body) {
    console.log(`Creating ${url} with ${JSON.stringify(body)}`)
    return fetch(url, {
        method: 'POST',
        body: JSON.stringify(body),
    })
    .then(r => r.json())
    .catch(err => {
        console.error(err)
        displayAlert('alert-danger', 'Create failed', false)
    });
}

function deleteUrl(url) {
    console.log(`Deleting ${url}`)
    return fetch(url, {
        method: 'DELETE'
    })
    .catch(err => {
        console.error(err)
        displayAlert('alert-danger', 'Delete failed', false)
    });
}

function displayAlert(type, text, autoDismiss) {
    let alertDiv = document.getElementById('alert');
    alertDiv.innerHTML = `
        <div class="alert ${type} alert-dismissible fade show w-100" role="alert">
          ${text}
          <button type="button" class="close" data-dismiss="alert" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
    `;
    if (autoDismiss) {
        setTimeout(function () {
            alertDiv.innerHTML = '';
        }, 2000);
    }
}

function createNumberInput(value) {
    return `<input type="number" class="form-control" placeholder="0.0" value="${value}" />`
}
