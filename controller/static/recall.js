
function updateAxis(preferences, dimension, prefix, value) {

    const left = document.getElementById(prefix + 'LeftId');
    const right = document.getElementById(prefix + 'RightId');

    if (Math.abs(value) < preferences.tolerance) {
        left.innerHTML = ''
        right.innerHTML = ''
        value = 0.0
    } else if (value < 0) {
        left.innerHTML = '&darr;'
        right.innerHTML = '&uarr;'
    } else {
        left.innerHTML = '&uarr;'
        right.innerHTML = '&darr;'
    }

    document.getElementById(prefix + 'ProgressId').innerHTML = createProgress(value)

    const distanceDiv = document.getElementById(prefix + 'DistanceId')
    const length = 2 * Math.pow(dimension, 2.0);
    const distance = Math.sqrt(length - (length * Math.cos(value * Math.PI / 180.0)));
    distanceDiv.innerHTML = `${distance.toPrecision(2)} ${preferences.dimensions.units}`
}

function updateLevel(preferences, position, current) {
    if (document.getElementById('recallId')) {
        updateAxis(preferences, preferences.dimensions.width, 'roll', current.roll - position.roll)
        updateAxis(preferences, preferences.dimensions.length, 'pitch', current.pitch - position.pitch)
        return true
    }
    return false
}

function constrainValue(value, min, max) {
    return Math.min(Math.max(value, min), max)
}

function calculateValueRatio(value, min, max) {
    const constrainedValue = constrainValue(value, min, max)
    return (constrainedValue - min) / (max - min)
}

function createProgress(angle) {
    const minValue = -4.0
    const maxValue = 4.0
    const tolerance = 0.1
    const width = 20.0
    const halfWidth = width / 2.0

    let relativeValue = 100.0 * calculateValueRatio(angle, minValue, maxValue)
    relativeValue = constrainValue(relativeValue, halfWidth, 100.0 - halfWidth)

    const lowerProgress = relativeValue - halfWidth
    const color = Math.abs(angle) < tolerance ? 'bg-success' : 'bg-primary'

    return `
        <div class="progress" style="height: 3rem;">
            <div class="progress-bar bg-light" role="progressbar" style="width: ${lowerProgress}%" aria-valuenow="1" aria-valuemin="0" aria-valuemax="100"></div>
            <div class="progress-bar ${color}" role="progressbar" style="width: ${width}%" aria-valuenow="1" aria-valuemin="0" aria-valuemax="100">${angle.toPrecision(2)}°</div>
        </div>
    `
}

function createGauge(prefix, name, leftName, rightName, units) {
    return `
        <div class="container text-light text-center p-0">
            <div class="row h3">
                <div class="col">${name}</div>
            </div>
            <div class="row no-gutters align-self-center h1">
                <div class="col" id="${prefix}LeftId"></div>
                <div class="col-8" id="${prefix}ProgressId">${createProgress(0)}</div>
                <div class="col" id="${prefix}RightId"></div>
            </div>
            <div class="row no-gutters h6">
                <div class="col">${leftName}</div>
                <div class="col-8 h3" id="${prefix}DistanceId">0.0 ${units}</div>
                <div class="col">${rightName}</div>
            </div>
        </div>
    `
}

function displayUpdate(timer, preferences, position) {
    fetch('api/corrected')
        .then(res => res.json())
        .then(current => {
            if (!updateLevel(preferences, position, current)) {
                clearInterval(timer)
            }
        })
        .catch(err => console.error(err));
}

async function displayLevel(matches, div) {
    const preferences = await getUrl('api/preference')
    displayRecall(div, {roll: 0.0, pitch: 0.0}, preferences, true)
}

async function displayPositionRecall(matches, div) {
    const json = await getUrls([
        'api/position/' + matches[1],
        'api/preference'
    ])
    displayRecall(div, json[0], json[1], false)
}

function displayRecall(div, position, preferences, isLevel) {
    div.innerHTML = `
        <div class="container h-100" id="recallId">
            <div class="row text-light mb-2">
                <div class="col text-center h3">
                    ${isLevel ? 'Level' : position.name}
                </div>
            </div>
            <div class="row h-25 align-items-center mb-2">
                <div class="col">
                    ${createGauge('roll', 'Roll', 'Left', 'Right', preferences.dimensions.units)}
                </div>
            </div>
            <div class="row h-25 align-items-center mb-3">
                <div class="col">
                    ${createGauge('pitch', 'Pitch', 'Front', 'Rear', preferences.dimensions.units)}
                </div>
            </div>
            <div class="row">
                <div class="col">
                    <a class="btn btn-lg btn-primary btn-block" role="button" href="#position/${position.id}" ${isLevel ? 'hidden' : ''}>Edit</a>
                </div>
            </div>
        </div>
    `;

    const timeout = Math.max(1000.0 / preferences.updateRate, 100)
    const timer = setInterval(() => {
        displayUpdate(timer, preferences, position)
    }, timeout)
}
