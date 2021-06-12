
function updateAxis(dimension, prefix, value) {
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
    distanceDiv.innerHTML = distance.toPrecision(2) + " " + preferences.dimensionUnits
}

function updateLevel(current) {
    updateAxis(preferences.dimensionWidth, 'roll', current.roll - (position ? position.roll : 0.0))
    updateAxis(preferences.dimensionLength, 'pitch', current.pitch - (position ? position.pitch : 0.0))
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

function displayUpdate() {
    getUrl('/api/corrected')
        .then(current => updateLevel(current))
        .catch(err => console.error(err));
}

document.getElementById('rollProgressId').innerHTML = createProgress(0)
document.getElementById('pitchProgressId').innerHTML = createProgress(0)

if (config.displaySseEnabled && (typeof (EventSource) !== "undefined")) {
    var source = new EventSource("/api/corrected/event");
    source.onmessage = function (event) {
        updateLevel(JSON.parse(event.data))
    };
} else {
    setInterval(
        () => displayUpdate(),
        1000.0 / config.displayUpdateRate
    )
}
