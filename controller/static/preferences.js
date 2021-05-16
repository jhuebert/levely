
async function displayPreferences(matches, div) {

    const preferences = await getUrl('api/preference')

    div.innerHTML = `
        <div class="container h-100">
            <div class="row px-3" id="alert"></div>
            <div class="text-light text-center h3 mb-2">Preferences</div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Width</label>
                <div class="col-sm-10" id="dimensionWidthId">
                    ${createNumberInput(preferences.dimensionWidth)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Length</label>
                <div class="col-sm-10" id="dimensionLengthId">
                    ${createNumberInput(preferences.dimensionLength)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Units</label>
                <div class="col-sm-10">
                    <select class="form-control" id="dimensionUnitsId" >
                        <option value="in" ${preferences.dimensionUnits === 'in' ? 'selected' : ''}>Inches</option>
                        <option value="cm" ${preferences.dimensionUnits === 'cm' ? 'selected' : ''}>Centimeters</option>
                    </select>
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Roll Axis</label>
                <div class="col-sm-10">
                    <select class="form-control">
                        <option id="orientationRollIdX"  value="x" ${preferences.orientationRoll === 'x' ? 'selected' : ''}>X</option>
                        <option id="orientationRollIdY"  value="Y" ${preferences.orientationRoll === 'y' ? 'selected' : ''}>Y</option>
                        <option id="orientationRollIdZ"  value="Z" ${preferences.orientationRoll === 'z' ? 'selected' : ''}>Z</option>
                    </select>
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Pitch Axis</label>
                <div class="col-sm-10">
                    <select class="form-control">
                        <option id="orientationPitchIdX"  value="x" ${preferences.orientationPitch === 'x' ? 'selected' : ''}>X</option>
                        <option id="orientationPitchIdY"  value="Y" ${preferences.orientationPitch === 'y' ? 'selected' : ''}>Y</option>
                        <option id="orientationPitchIdZ"  value="Z" ${preferences.orientationPitch === 'z' ? 'selected' : ''}>Z</option>
                    </select>
                </div>
            </div>
            <div class="custom-control custom-switch">
              <input type="checkbox" class="custom-control-input" id="orientationInvertRollId" ${preferences.orientationInvertRoll ? 'checked' : ''}>
              <label class="custom-control-label text-light h5" for="orientationInvertRollId">Invert Roll</label>
            </div>
            <div class="custom-control custom-switch">
              <input type="checkbox" class="custom-control-input" id="orientationInvertPitchId" ${preferences.orientationInvertPitch ? 'checked' : ''}>
              <label class="custom-control-label text-light h5" for="orientationInvertPitchId">Invert Pitch</label>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Level Tolerance</label>
                <div class="col-sm-10" id="levelToleranceId">
                    ${createNumberInput(preferences.levelTolerance)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Display Rate</label>
                <div class="col-sm-10" id="displayRateId">
                    ${createNumberInput(preferences.displayRate)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Accelerometer Rate</label>
                <div class="col-sm-10" id="accelerometerRateId">
                    ${createNumberInput(preferences.accelerometerRate)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Accelerometer Smoothing</label>
                <div class="col-sm-10" id="accelerometerSmoothingId">
                    ${createNumberInput(preferences.accelerometerSmoothing)}
                </div>
            </div>
            <button type="button" class="btn btn-lg btn-primary btn-block mt-5" onclick="savePreferences()">Save</button>
            <!-- TODO Add button to export the database -->
        </div>
    `;
}

async function savePreferences() {

    const preferences = {}
    preferences.dimensionWidth = parseFloat(document.getElementById('dimensionWidthId').firstElementChild.value)
    preferences.dimensionLength = parseFloat(document.getElementById('dimensionLengthId').firstElementChild.value)
    preferences.dimensionUnits = document.getElementById('dimensionUnitsId').value

    preferences.orientationRoll = document.getElementById('orientationRollIdX').selected ? 'x' : preferences.orientationRoll
    preferences.orientationRoll = document.getElementById('orientationRollIdY').selected ? 'y' : preferences.orientationRoll
    preferences.orientationRoll = document.getElementById('orientationRollIdZ').selected ? 'z' : preferences.orientationRoll

    preferences.orientationPitch = document.getElementById('orientationPitchIdX').selected ? 'x' : preferences.orientationPitch
    preferences.orientationPitch = document.getElementById('orientationPitchIdY').selected ? 'y' : preferences.orientationPitch
    preferences.orientationPitch = document.getElementById('orientationPitchIdZ').selected ? 'z' : preferences.orientationPitch

    preferences.orientationInvertRoll = document.getElementById('orientationInvertRollId').checked
    preferences.orientationInvertPitch = document.getElementById('orientationInvertPitchId').checked

    preferences.levelTolerance = parseFloat(document.getElementById('levelToleranceId').firstElementChild.value)
    preferences.displayRate = parseFloat(document.getElementById('displayRateId').firstElementChild.value)
    preferences.accelerometerRate = parseFloat(document.getElementById('accelerometerRateId').firstElementChild.value)
    preferences.accelerometerSmoothing = parseFloat(document.getElementById('accelerometerSmoothingId').firstElementChild.value)

    console.log(preferences)

    const validationResult = validatePreferences(preferences)
    if (validationResult) {
        displayAlert('alert-danger', validationResult)
        return
    }

    await putUrl('api/preference', preferences)

    displayAlert('alert-success', 'Update success', true)
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

    if (!preferences.levelTolerance || (preferences.levelTolerance <= 0.0)) {
        return 'Level tolerance must be greater than zero'
    }

    if (!preferences.accelerometerRate || (preferences.accelerometerRate <= 0.0)) {
        return 'Accelerometer rate must be greater than zero'
    }

    if (!preferences.accelerometerSmoothing || (preferences.accelerometerSmoothing <= 0.0)) {
        return 'Accelerometer rate must be greater than zero'
    }

    return null
}
