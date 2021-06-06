
async function savePreferences() {
    const preferences = {}
    preferences.dimensionWidth = parseFloat(document.getElementById('dimensionWidthId').value)
    preferences.dimensionLength = parseFloat(document.getElementById('dimensionLengthId').value)
    preferences.dimensionUnits = document.getElementById('dimensionUnitsId').value
    preferences.orientationRoll = document.getElementById('orientationRollId').value
    preferences.orientationPitch = document.getElementById('orientationPitchId').value
    preferences.orientationInvertRoll = document.getElementById('orientationInvertRollId').checked
    preferences.orientationInvertPitch = document.getElementById('orientationInvertPitchId').checked
    const validationResult = validatePreferences(preferences)
    if (validationResult) {
        displayAlert('alert-danger', validationResult)
        return
    }
    await putUrl('/api/preference', preferences)
}

function validatePreferences(preferences) {
    if (!preferences.dimensionWidth || (preferences.dimensionWidth <= 0.0)) {
        return 'Width must be greater than zero'
    }
    if (!preferences.dimensionLength || (preferences.dimensionLength <= 0.0)) {
        return 'Length must be greater than zero'
    }
    if (!preferences.dimensionUnits) {
        return 'Dimension unit not specified'
    }
    if (!preferences.orientationRoll) {
        return "Roll axis is not set"
    }
    if (!preferences.orientationPitch) {
        return "Pitch axis is not set"
    }
    if (preferences.orientationRoll === preferences.orientationPitch) {
        return "Pitch and roll axis must not be the same value"
    }
    return null
}
