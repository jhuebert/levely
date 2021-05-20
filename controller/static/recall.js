
function updateAxis(preferences, config, dimension, prefix, value) {

    const left = document.getElementById(prefix + 'LeftId');
    const right = document.getElementById(prefix + 'RightId');

    if (Math.abs(value) < config.displayLevelTolerance) {
        left.innerHTML = ''
        right.innerHTML = ''
        value = 0.0
    } else if (value < 0) {
        left.innerHTML = '&#9660;'
        right.innerHTML = '&#9650;'
    } else {
        left.innerHTML = '&#9650;'
        right.innerHTML = '&#9660;'
    }

    document.getElementById(prefix + 'ProgressId').innerHTML = createProgress(value)

    const distanceDiv = document.getElementById(prefix + 'DistanceId')
    const length = 2 * Math.pow(dimension, 2.0);
    const distance = Math.sqrt(length - (length * Math.cos(value * Math.PI / 180.0)));
    distanceDiv.innerHTML = `${distance.toPrecision(2)} ${preferences.dimensionUnits}`
}

function updateLevel(preferences, config, position, current) {
    if (document.getElementById('recallId')) {
        updateAxis(preferences, config, preferences.dimensionWidth, 'roll', current.roll - position.roll)
        updateAxis(preferences, config, preferences.dimensionLength, 'pitch', current.pitch - position.pitch)
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

function displayUpdate(timer, preferences, config, position) {
    fetch('api/corrected')
        .then(res => res.json())
        .then(current => {
            if (!updateLevel(preferences, config, position, current)) {
                clearInterval(timer)
            }
        })
        .catch(err => console.error(err));
}

async function displayLevel(matches, div) {
    const json = await getUrls([
        'api/preference',
        'api/config'
    ])
    displayRecall(div, {roll: 0.0, pitch: 0.0}, json[0], json[1], true)
}

async function displayPositionRecall(matches, div) {
    const json = await getUrls([
        'api/position/' + matches[1],
        'api/preference',
        'api/config'
    ])
    displayRecall(div, json[0], json[1], json[2], false)
}

function displayRecall(div, position, preferences, config, isLevel) {
    div.innerHTML = `
        <div class="container h-100" id="recallId">
            <div class="row text-light mb-2">
                <div class="col text-center h3">
                    ${isLevel ? 'Level' : position.name}
                </div>
            </div>
            <div class="row h-25 align-items-center mb-2">
                <div class="col">
                    ${createGauge('roll', 'Roll', 'Left', 'Right', preferences.dimensionUnits)}
                </div>
            </div>
            <div class="row h-25 align-items-center mb-3">
                <div class="col">
                    ${createGauge('pitch', 'Pitch', 'Front', 'Rear', preferences.dimensionUnits)}
                </div>
            </div>
            <div class="row">
                <div class="col">
                    <a class="btn btn-lg btn-primary btn-block" role="button" href="#position/${position.id}" ${isLevel ? 'hidden' : ''}>Edit</a>
                </div>
            </div>
        </div>
    `;

    // Limit to a maximum of 10Hz
    const timeout = Math.max(1000.0 / config.displayUpdateRate, 100)
    const timer = setInterval(() => {
        displayUpdate(timer, preferences, config, position)
    }, timeout)
}
